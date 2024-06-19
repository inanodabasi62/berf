package app

import (
	"io"
	"os"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/config"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/log"
	tmdb "github.com/tendermint/tm-db"
)

const appName = "berf"

var (
	DefaultNodeHome = os.ExpandEnv("$HOME/.berf")
	maccPerms       = map[string][]string{
		authtypes.FeeCollectorName:     {authtypes.Minter, authtypes.Burner},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
	}
)

type App struct {
	*baseapp.BaseApp
	cdc               *codec.LegacyAmino
	appCodec          codec.Codec
	interfaceRegistry types.InterfaceRegistry

	keys    map[string]*sdk.KVStoreKey
	tkeys   map[string]*sdk.TransientStoreKey
	memKeys map[string]*sdk.MemoryStoreKey

	AccountKeeper authkeeper.AccountKeeper
	BankKeeper    bankkeeper.Keeper
	StakingKeeper stakingkeeper.Keeper

	mm *module.Manager
}

func NewRootCmd() (*cobra.Command, params.EncodingConfig) {
	encodingConfig := MakeEncodingConfig()
	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastBlock).
		WithHomeDir(DefaultNodeHome).
		WithViper("")

	rootCmd := &cobra.Command{
		Use:   appName,
		Short: "Start berf node",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.OutOrStderr())

			initClientCtx, err := client.ReadPersistentCommandFlags(initClientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			initClientCtx, err = config.ReadFromClientConfig(initClientCtx)
			if err != nil {
				return err
			}

			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			tmCfg := sdk.GetConfig()
			tmCfg.Seal()

			return server.InterceptConfigsPreRunHandler(cmd, "", nil)
		},
	}

	server.AddCommands(rootCmd, DefaultNodeHome, newApp, createSimappAndExport, addModuleInitFlags)
	rootCmd.AddCommand(rpc.StatusCommand())
	rootCmd.AddCommand(genutilcli.InitCmd(DefaultNodeHome, appName))
	rootCmd.AddCommand(genutilcli.CollectGenTxsCmd(banktypes.GenesisBalancesIterator{}, DefaultNodeHome))
	rootCmd.AddCommand(genutilcli.GenTxCmd(banktypes.GenesisBalancesIterator{}, DefaultNodeHome))
	rootCmd.AddCommand(genutilcli.ValidateGenesisCmd())
	rootCmd.AddCommand(client.NewCompletionCmd(rootCmd, true))

	return rootCmd, encodingConfig
}

func newApp(logger log.Logger, db tmdb.DB, traceStore io.Writer, appOpts server.AppOptions) *App {
	cdc := MakeEncodingConfig().Marshaler

	bApp := baseapp.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseapp.SetPruning(store.NewPruningOptionsFromString(viper.GetString("pruning"))))
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)

	keys := sdk.NewKVStoreKeys(
		authtypes.StoreKey, stakingtypes.StoreKey, banktypes.StoreKey, genutiltypes.StoreKey,
	)
	tkeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey)
	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	app := &App{
		BaseApp:           bApp,
		cdc:               cdc,
		appCodec:          cdc,
		interfaceRegistry: encodingConfig.InterfaceRegistry,
		keys:              keys,
		tkeys:             tkeys,
		memKeys:           memKeys,
	}

	app.ParamsKeeper = initParamsKeeper(appCodec, legacyAmino, keys[paramstypes.StoreKey], tkeys[paramstypes.TStoreKey])
	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec, keys[authtypes.StoreKey], app.GetSubspace(authtypes.ModuleName), authtypes.ProtoBaseAccount, maccPerms,
	)
	app.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec, keys[banktypes.StoreKey], app.AccountKeeper, app.GetSubspace(banktypes.ModuleName), app.ModuleAccountAddrs(),
	)
	stakingKeeper := stakingkeeper.NewKeeper(
		appCodec, keys[stakingtypes.StoreKey], app.AccountKeeper, app.BankKeeper, app.GetSubspace(stakingtypes.ModuleName),
	)
	app.StakingKeeper = *stakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(),
	)

	app.mm = module.NewManager(
		genutil.NewAppModule(app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(appCodec, app.AccountKeeper, nil),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper),
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
	)

	app.mm.SetOrderBeginBlockers(upgradetypes.ModuleName, capabilitytypes.ModuleName, genutiltypes.ModuleName)
	app.mm.SetOrderEndBlockers(stakingtypes.ModuleName)

	app.mm.RegisterInvariants(app.CrisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter(), codec.NewAminoCodec(cdc))

	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	return app
}

func MakeEncodingConfig() params.EncodingConfig {
	encodingConfig := params.MakeEncodingConfig()
	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}

func createSimappAndExport(logger log.Logger, db tmdb.DB, traceStore io.Writer, appOpts server.AppOptions) (server.ExportedApp, error) {
	var simApp *App
	name := appName

	homePath := viper.GetString(client.HomeFlag)
	loadLatest := viper.GetBool(server.FlagOverwriteHeight)

	if homePath == "" {
		homePath = DefaultNodeHome
	}

	simApp = newApp(logger, db, traceStore, appOpts)

	return simApp.ExportAppStateAndValidators(loadLatest, []string{})
}

func addModuleInitFlags(startCmd *cobra.Command) {
	// add any additional init command flags
}
