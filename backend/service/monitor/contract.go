// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package monitor

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// MonitorMetaData contains all meta data concerning the Monitor contract.
var MonitorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount0In\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount1In\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount0Out\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount1Out\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"Buy\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"totalSupply\",\"type\":\"uint256\"}],\"name\":\"CreateToken\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"Launch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"reserve0\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"reserve1\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"k\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"}],\"name\":\"Liquity\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"ReadyToLaunch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountOIn\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount1In\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount0Out\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount1Out\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"Sell\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MANAGE_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"OPERATE_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"buy\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"}],\"name\":\"createToken\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenInAmount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"}],\"name\":\"getAmountOut\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"getRoleMember\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleMemberCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIRouter\",\"name\":\"_iRouter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_feeReceiver\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_operator\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"launch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"sellAmount\",\"type\":\"uint256\"}],\"name\":\"sell\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
}

// MonitorABI is the input ABI used to generate the binding from.
// Deprecated: Use MonitorMetaData.ABI instead.
var MonitorABI = MonitorMetaData.ABI

// Monitor is an auto generated Go binding around an Ethereum contract.
type Monitor struct {
	MonitorCaller     // Read-only binding to the contract
	MonitorTransactor // Write-only binding to the contract
	MonitorFilterer   // Log filterer for contract events
}

// MonitorCaller is an auto generated read-only Go binding around an Ethereum contract.
type MonitorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MonitorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MonitorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MonitorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MonitorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MonitorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MonitorSession struct {
	Contract     *Monitor          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MonitorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MonitorCallerSession struct {
	Contract *MonitorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// MonitorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MonitorTransactorSession struct {
	Contract     *MonitorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// MonitorRaw is an auto generated low-level Go binding around an Ethereum contract.
type MonitorRaw struct {
	Contract *Monitor // Generic contract binding to access the raw methods on
}

// MonitorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MonitorCallerRaw struct {
	Contract *MonitorCaller // Generic read-only contract binding to access the raw methods on
}

// MonitorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MonitorTransactorRaw struct {
	Contract *MonitorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMonitor creates a new instance of Monitor, bound to a specific deployed contract.
func NewMonitor(address common.Address, backend bind.ContractBackend) (*Monitor, error) {
	contract, err := bindMonitor(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Monitor{MonitorCaller: MonitorCaller{contract: contract}, MonitorTransactor: MonitorTransactor{contract: contract}, MonitorFilterer: MonitorFilterer{contract: contract}}, nil
}

// NewMonitorCaller creates a new read-only instance of Monitor, bound to a specific deployed contract.
func NewMonitorCaller(address common.Address, caller bind.ContractCaller) (*MonitorCaller, error) {
	contract, err := bindMonitor(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MonitorCaller{contract: contract}, nil
}

// NewMonitorTransactor creates a new write-only instance of Monitor, bound to a specific deployed contract.
func NewMonitorTransactor(address common.Address, transactor bind.ContractTransactor) (*MonitorTransactor, error) {
	contract, err := bindMonitor(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MonitorTransactor{contract: contract}, nil
}

// NewMonitorFilterer creates a new log filterer instance of Monitor, bound to a specific deployed contract.
func NewMonitorFilterer(address common.Address, filterer bind.ContractFilterer) (*MonitorFilterer, error) {
	contract, err := bindMonitor(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MonitorFilterer{contract: contract}, nil
}

// bindMonitor binds a generic wrapper to an already deployed contract.
func bindMonitor(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MonitorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Monitor *MonitorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Monitor.Contract.MonitorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Monitor *MonitorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Monitor.Contract.MonitorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Monitor *MonitorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Monitor.Contract.MonitorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Monitor *MonitorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Monitor.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Monitor *MonitorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Monitor.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Monitor *MonitorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Monitor.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Monitor *MonitorCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Monitor.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Monitor *MonitorSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Monitor.Contract.DEFAULTADMINROLE(&_Monitor.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Monitor *MonitorCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Monitor.Contract.DEFAULTADMINROLE(&_Monitor.CallOpts)
}

// MANAGEROLE is a free data retrieval call binding the contract method 0x60a4b76a.
//
// Solidity: function MANAGE_ROLE() view returns(bytes32)
func (_Monitor *MonitorCaller) MANAGEROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Monitor.contract.Call(opts, &out, "MANAGE_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MANAGEROLE is a free data retrieval call binding the contract method 0x60a4b76a.
//
// Solidity: function MANAGE_ROLE() view returns(bytes32)
func (_Monitor *MonitorSession) MANAGEROLE() ([32]byte, error) {
	return _Monitor.Contract.MANAGEROLE(&_Monitor.CallOpts)
}

// MANAGEROLE is a free data retrieval call binding the contract method 0x60a4b76a.
//
// Solidity: function MANAGE_ROLE() view returns(bytes32)
func (_Monitor *MonitorCallerSession) MANAGEROLE() ([32]byte, error) {
	return _Monitor.Contract.MANAGEROLE(&_Monitor.CallOpts)
}

// OPERATEROLE is a free data retrieval call binding the contract method 0xc81f0af8.
//
// Solidity: function OPERATE_ROLE() view returns(bytes32)
func (_Monitor *MonitorCaller) OPERATEROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Monitor.contract.Call(opts, &out, "OPERATE_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// OPERATEROLE is a free data retrieval call binding the contract method 0xc81f0af8.
//
// Solidity: function OPERATE_ROLE() view returns(bytes32)
func (_Monitor *MonitorSession) OPERATEROLE() ([32]byte, error) {
	return _Monitor.Contract.OPERATEROLE(&_Monitor.CallOpts)
}

// OPERATEROLE is a free data retrieval call binding the contract method 0xc81f0af8.
//
// Solidity: function OPERATE_ROLE() view returns(bytes32)
func (_Monitor *MonitorCallerSession) OPERATEROLE() ([32]byte, error) {
	return _Monitor.Contract.OPERATEROLE(&_Monitor.CallOpts)
}

// GetAmountOut is a free data retrieval call binding the contract method 0xff9c8ac6.
//
// Solidity: function getAmountOut(address tokenIn, uint256 tokenInAmount, address tokenOut) view returns(uint256)
func (_Monitor *MonitorCaller) GetAmountOut(opts *bind.CallOpts, tokenIn common.Address, tokenInAmount *big.Int, tokenOut common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Monitor.contract.Call(opts, &out, "getAmountOut", tokenIn, tokenInAmount, tokenOut)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAmountOut is a free data retrieval call binding the contract method 0xff9c8ac6.
//
// Solidity: function getAmountOut(address tokenIn, uint256 tokenInAmount, address tokenOut) view returns(uint256)
func (_Monitor *MonitorSession) GetAmountOut(tokenIn common.Address, tokenInAmount *big.Int, tokenOut common.Address) (*big.Int, error) {
	return _Monitor.Contract.GetAmountOut(&_Monitor.CallOpts, tokenIn, tokenInAmount, tokenOut)
}

// GetAmountOut is a free data retrieval call binding the contract method 0xff9c8ac6.
//
// Solidity: function getAmountOut(address tokenIn, uint256 tokenInAmount, address tokenOut) view returns(uint256)
func (_Monitor *MonitorCallerSession) GetAmountOut(tokenIn common.Address, tokenInAmount *big.Int, tokenOut common.Address) (*big.Int, error) {
	return _Monitor.Contract.GetAmountOut(&_Monitor.CallOpts, tokenIn, tokenInAmount, tokenOut)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Monitor *MonitorCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _Monitor.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Monitor *MonitorSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Monitor.Contract.GetRoleAdmin(&_Monitor.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Monitor *MonitorCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Monitor.Contract.GetRoleAdmin(&_Monitor.CallOpts, role)
}

// GetRoleMember is a free data retrieval call binding the contract method 0x9010d07c.
//
// Solidity: function getRoleMember(bytes32 role, uint256 index) view returns(address)
func (_Monitor *MonitorCaller) GetRoleMember(opts *bind.CallOpts, role [32]byte, index *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Monitor.contract.Call(opts, &out, "getRoleMember", role, index)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetRoleMember is a free data retrieval call binding the contract method 0x9010d07c.
//
// Solidity: function getRoleMember(bytes32 role, uint256 index) view returns(address)
func (_Monitor *MonitorSession) GetRoleMember(role [32]byte, index *big.Int) (common.Address, error) {
	return _Monitor.Contract.GetRoleMember(&_Monitor.CallOpts, role, index)
}

// GetRoleMember is a free data retrieval call binding the contract method 0x9010d07c.
//
// Solidity: function getRoleMember(bytes32 role, uint256 index) view returns(address)
func (_Monitor *MonitorCallerSession) GetRoleMember(role [32]byte, index *big.Int) (common.Address, error) {
	return _Monitor.Contract.GetRoleMember(&_Monitor.CallOpts, role, index)
}

// GetRoleMemberCount is a free data retrieval call binding the contract method 0xca15c873.
//
// Solidity: function getRoleMemberCount(bytes32 role) view returns(uint256)
func (_Monitor *MonitorCaller) GetRoleMemberCount(opts *bind.CallOpts, role [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _Monitor.contract.Call(opts, &out, "getRoleMemberCount", role)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRoleMemberCount is a free data retrieval call binding the contract method 0xca15c873.
//
// Solidity: function getRoleMemberCount(bytes32 role) view returns(uint256)
func (_Monitor *MonitorSession) GetRoleMemberCount(role [32]byte) (*big.Int, error) {
	return _Monitor.Contract.GetRoleMemberCount(&_Monitor.CallOpts, role)
}

// GetRoleMemberCount is a free data retrieval call binding the contract method 0xca15c873.
//
// Solidity: function getRoleMemberCount(bytes32 role) view returns(uint256)
func (_Monitor *MonitorCallerSession) GetRoleMemberCount(role [32]byte) (*big.Int, error) {
	return _Monitor.Contract.GetRoleMemberCount(&_Monitor.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Monitor *MonitorCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _Monitor.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Monitor *MonitorSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Monitor.Contract.HasRole(&_Monitor.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Monitor *MonitorCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Monitor.Contract.HasRole(&_Monitor.CallOpts, role, account)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Monitor *MonitorCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Monitor.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Monitor *MonitorSession) ProxiableUUID() ([32]byte, error) {
	return _Monitor.Contract.ProxiableUUID(&_Monitor.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Monitor *MonitorCallerSession) ProxiableUUID() ([32]byte, error) {
	return _Monitor.Contract.ProxiableUUID(&_Monitor.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Monitor *MonitorCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _Monitor.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Monitor *MonitorSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Monitor.Contract.SupportsInterface(&_Monitor.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Monitor *MonitorCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Monitor.Contract.SupportsInterface(&_Monitor.CallOpts, interfaceId)
}

// Buy is a paid mutator transaction binding the contract method 0xf088d547.
//
// Solidity: function buy(address token) payable returns()
func (_Monitor *MonitorTransactor) Buy(opts *bind.TransactOpts, token common.Address) (*types.Transaction, error) {
	return _Monitor.contract.Transact(opts, "buy", token)
}

// Buy is a paid mutator transaction binding the contract method 0xf088d547.
//
// Solidity: function buy(address token) payable returns()
func (_Monitor *MonitorSession) Buy(token common.Address) (*types.Transaction, error) {
	return _Monitor.Contract.Buy(&_Monitor.TransactOpts, token)
}

// Buy is a paid mutator transaction binding the contract method 0xf088d547.
//
// Solidity: function buy(address token) payable returns()
func (_Monitor *MonitorTransactorSession) Buy(token common.Address) (*types.Transaction, error) {
	return _Monitor.Contract.Buy(&_Monitor.TransactOpts, token)
}

// CreateToken is a paid mutator transaction binding the contract method 0xb1037e78.
//
// Solidity: function createToken(uint256 tokenId, string name, string symbol) payable returns()
func (_Monitor *MonitorTransactor) CreateToken(opts *bind.TransactOpts, tokenId *big.Int, name string, symbol string) (*types.Transaction, error) {
	return _Monitor.contract.Transact(opts, "createToken", tokenId, name, symbol)
}

// CreateToken is a paid mutator transaction binding the contract method 0xb1037e78.
//
// Solidity: function createToken(uint256 tokenId, string name, string symbol) payable returns()
func (_Monitor *MonitorSession) CreateToken(tokenId *big.Int, name string, symbol string) (*types.Transaction, error) {
	return _Monitor.Contract.CreateToken(&_Monitor.TransactOpts, tokenId, name, symbol)
}

// CreateToken is a paid mutator transaction binding the contract method 0xb1037e78.
//
// Solidity: function createToken(uint256 tokenId, string name, string symbol) payable returns()
func (_Monitor *MonitorTransactorSession) CreateToken(tokenId *big.Int, name string, symbol string) (*types.Transaction, error) {
	return _Monitor.Contract.CreateToken(&_Monitor.TransactOpts, tokenId, name, symbol)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Monitor *MonitorTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Monitor.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Monitor *MonitorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Monitor.Contract.GrantRole(&_Monitor.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Monitor *MonitorTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Monitor.Contract.GrantRole(&_Monitor.TransactOpts, role, account)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address _iRouter, address _feeReceiver, address _operator) returns()
func (_Monitor *MonitorTransactor) Initialize(opts *bind.TransactOpts, _iRouter common.Address, _feeReceiver common.Address, _operator common.Address) (*types.Transaction, error) {
	return _Monitor.contract.Transact(opts, "initialize", _iRouter, _feeReceiver, _operator)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address _iRouter, address _feeReceiver, address _operator) returns()
func (_Monitor *MonitorSession) Initialize(_iRouter common.Address, _feeReceiver common.Address, _operator common.Address) (*types.Transaction, error) {
	return _Monitor.Contract.Initialize(&_Monitor.TransactOpts, _iRouter, _feeReceiver, _operator)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address _iRouter, address _feeReceiver, address _operator) returns()
func (_Monitor *MonitorTransactorSession) Initialize(_iRouter common.Address, _feeReceiver common.Address, _operator common.Address) (*types.Transaction, error) {
	return _Monitor.Contract.Initialize(&_Monitor.TransactOpts, _iRouter, _feeReceiver, _operator)
}

// Launch is a paid mutator transaction binding the contract method 0x214013ca.
//
// Solidity: function launch(address token) returns()
func (_Monitor *MonitorTransactor) Launch(opts *bind.TransactOpts, token common.Address) (*types.Transaction, error) {
	return _Monitor.contract.Transact(opts, "launch", token)
}

// Launch is a paid mutator transaction binding the contract method 0x214013ca.
//
// Solidity: function launch(address token) returns()
func (_Monitor *MonitorSession) Launch(token common.Address) (*types.Transaction, error) {
	return _Monitor.Contract.Launch(&_Monitor.TransactOpts, token)
}

// Launch is a paid mutator transaction binding the contract method 0x214013ca.
//
// Solidity: function launch(address token) returns()
func (_Monitor *MonitorTransactorSession) Launch(token common.Address) (*types.Transaction, error) {
	return _Monitor.Contract.Launch(&_Monitor.TransactOpts, token)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Monitor *MonitorTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Monitor.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Monitor *MonitorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Monitor.Contract.RenounceRole(&_Monitor.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Monitor *MonitorTransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Monitor.Contract.RenounceRole(&_Monitor.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Monitor *MonitorTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Monitor.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Monitor *MonitorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Monitor.Contract.RevokeRole(&_Monitor.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Monitor *MonitorTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Monitor.Contract.RevokeRole(&_Monitor.TransactOpts, role, account)
}

// Sell is a paid mutator transaction binding the contract method 0x6c197ff5.
//
// Solidity: function sell(address token, uint256 sellAmount) returns()
func (_Monitor *MonitorTransactor) Sell(opts *bind.TransactOpts, token common.Address, sellAmount *big.Int) (*types.Transaction, error) {
	return _Monitor.contract.Transact(opts, "sell", token, sellAmount)
}

// Sell is a paid mutator transaction binding the contract method 0x6c197ff5.
//
// Solidity: function sell(address token, uint256 sellAmount) returns()
func (_Monitor *MonitorSession) Sell(token common.Address, sellAmount *big.Int) (*types.Transaction, error) {
	return _Monitor.Contract.Sell(&_Monitor.TransactOpts, token, sellAmount)
}

// Sell is a paid mutator transaction binding the contract method 0x6c197ff5.
//
// Solidity: function sell(address token, uint256 sellAmount) returns()
func (_Monitor *MonitorTransactorSession) Sell(token common.Address, sellAmount *big.Int) (*types.Transaction, error) {
	return _Monitor.Contract.Sell(&_Monitor.TransactOpts, token, sellAmount)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_Monitor *MonitorTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _Monitor.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_Monitor *MonitorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _Monitor.Contract.UpgradeTo(&_Monitor.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_Monitor *MonitorTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _Monitor.Contract.UpgradeTo(&_Monitor.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Monitor *MonitorTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Monitor.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Monitor *MonitorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Monitor.Contract.UpgradeToAndCall(&_Monitor.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Monitor *MonitorTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Monitor.Contract.UpgradeToAndCall(&_Monitor.TransactOpts, newImplementation, data)
}

// MonitorAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the Monitor contract.
type MonitorAdminChangedIterator struct {
	Event *MonitorAdminChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MonitorAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MonitorAdminChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MonitorAdminChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MonitorAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MonitorAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MonitorAdminChanged represents a AdminChanged event raised by the Monitor contract.
type MonitorAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_Monitor *MonitorFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*MonitorAdminChangedIterator, error) {

	logs, sub, err := _Monitor.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &MonitorAdminChangedIterator{contract: _Monitor.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_Monitor *MonitorFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *MonitorAdminChanged) (event.Subscription, error) {

	logs, sub, err := _Monitor.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MonitorAdminChanged)
				if err := _Monitor.contract.UnpackLog(event, "AdminChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAdminChanged is a log parse operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_Monitor *MonitorFilterer) ParseAdminChanged(log types.Log) (*MonitorAdminChanged, error) {
	event := new(MonitorAdminChanged)
	if err := _Monitor.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MonitorBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the Monitor contract.
type MonitorBeaconUpgradedIterator struct {
	Event *MonitorBeaconUpgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MonitorBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MonitorBeaconUpgraded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MonitorBeaconUpgraded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MonitorBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MonitorBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MonitorBeaconUpgraded represents a BeaconUpgraded event raised by the Monitor contract.
type MonitorBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_Monitor *MonitorFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*MonitorBeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _Monitor.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &MonitorBeaconUpgradedIterator{contract: _Monitor.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_Monitor *MonitorFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *MonitorBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _Monitor.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MonitorBeaconUpgraded)
				if err := _Monitor.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBeaconUpgraded is a log parse operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_Monitor *MonitorFilterer) ParseBeaconUpgraded(log types.Log) (*MonitorBeaconUpgraded, error) {
	event := new(MonitorBeaconUpgraded)
	if err := _Monitor.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MonitorBuyIterator is returned from FilterBuy and is used to iterate over the raw logs and unpacked data for Buy events raised by the Monitor contract.
type MonitorBuyIterator struct {
	Event *MonitorBuy // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MonitorBuyIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MonitorBuy)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MonitorBuy)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MonitorBuyIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MonitorBuyIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MonitorBuy represents a Buy event raised by the Monitor contract.
type MonitorBuy struct {
	Sender     common.Address
	Token      common.Address
	Amount0In  *big.Int
	Amount1In  *big.Int
	Amount0Out *big.Int
	Amount1Out *big.Int
	To         common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterBuy is a free log retrieval operation binding the contract event 0x8d0aaf9f88a0c8ece371b7250c39e5b2f8ea6338392a3cdd636d0800e24f382c.
//
// Solidity: event Buy(address sender, address token, uint256 amount0In, uint256 amount1In, uint256 amount0Out, uint256 amount1Out, address to)
func (_Monitor *MonitorFilterer) FilterBuy(opts *bind.FilterOpts) (*MonitorBuyIterator, error) {

	logs, sub, err := _Monitor.contract.FilterLogs(opts, "Buy")
	if err != nil {
		return nil, err
	}
	return &MonitorBuyIterator{contract: _Monitor.contract, event: "Buy", logs: logs, sub: sub}, nil
}

// WatchBuy is a free log subscription operation binding the contract event 0x8d0aaf9f88a0c8ece371b7250c39e5b2f8ea6338392a3cdd636d0800e24f382c.
//
// Solidity: event Buy(address sender, address token, uint256 amount0In, uint256 amount1In, uint256 amount0Out, uint256 amount1Out, address to)
func (_Monitor *MonitorFilterer) WatchBuy(opts *bind.WatchOpts, sink chan<- *MonitorBuy) (event.Subscription, error) {

	logs, sub, err := _Monitor.contract.WatchLogs(opts, "Buy")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MonitorBuy)
				if err := _Monitor.contract.UnpackLog(event, "Buy", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBuy is a log parse operation binding the contract event 0x8d0aaf9f88a0c8ece371b7250c39e5b2f8ea6338392a3cdd636d0800e24f382c.
//
// Solidity: event Buy(address sender, address token, uint256 amount0In, uint256 amount1In, uint256 amount0Out, uint256 amount1Out, address to)
func (_Monitor *MonitorFilterer) ParseBuy(log types.Log) (*MonitorBuy, error) {
	event := new(MonitorBuy)
	if err := _Monitor.contract.UnpackLog(event, "Buy", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MonitorCreateTokenIterator is returned from FilterCreateToken and is used to iterate over the raw logs and unpacked data for CreateToken events raised by the Monitor contract.
type MonitorCreateTokenIterator struct {
	Event *MonitorCreateToken // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MonitorCreateTokenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MonitorCreateToken)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MonitorCreateToken)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MonitorCreateTokenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MonitorCreateTokenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MonitorCreateToken represents a CreateToken event raised by the Monitor contract.
type MonitorCreateToken struct {
	TokenId     *big.Int
	Token       common.Address
	Name        string
	Symbol      string
	TotalSupply *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterCreateToken is a free log retrieval operation binding the contract event 0x1c56be39f0a7cd9b5f7d349a32369c4514e78870028ea4b721570f9a0ea3b127.
//
// Solidity: event CreateToken(uint256 tokenId, address token, string name, string symbol, uint256 totalSupply)
func (_Monitor *MonitorFilterer) FilterCreateToken(opts *bind.FilterOpts) (*MonitorCreateTokenIterator, error) {

	logs, sub, err := _Monitor.contract.FilterLogs(opts, "CreateToken")
	if err != nil {
		return nil, err
	}
	return &MonitorCreateTokenIterator{contract: _Monitor.contract, event: "CreateToken", logs: logs, sub: sub}, nil
}

// WatchCreateToken is a free log subscription operation binding the contract event 0x1c56be39f0a7cd9b5f7d349a32369c4514e78870028ea4b721570f9a0ea3b127.
//
// Solidity: event CreateToken(uint256 tokenId, address token, string name, string symbol, uint256 totalSupply)
func (_Monitor *MonitorFilterer) WatchCreateToken(opts *bind.WatchOpts, sink chan<- *MonitorCreateToken) (event.Subscription, error) {

	logs, sub, err := _Monitor.contract.WatchLogs(opts, "CreateToken")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MonitorCreateToken)
				if err := _Monitor.contract.UnpackLog(event, "CreateToken", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCreateToken is a log parse operation binding the contract event 0x1c56be39f0a7cd9b5f7d349a32369c4514e78870028ea4b721570f9a0ea3b127.
//
// Solidity: event CreateToken(uint256 tokenId, address token, string name, string symbol, uint256 totalSupply)
func (_Monitor *MonitorFilterer) ParseCreateToken(log types.Log) (*MonitorCreateToken, error) {
	event := new(MonitorCreateToken)
	if err := _Monitor.contract.UnpackLog(event, "CreateToken", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MonitorInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Monitor contract.
type MonitorInitializedIterator struct {
	Event *MonitorInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MonitorInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MonitorInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MonitorInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MonitorInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MonitorInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MonitorInitialized represents a Initialized event raised by the Monitor contract.
type MonitorInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Monitor *MonitorFilterer) FilterInitialized(opts *bind.FilterOpts) (*MonitorInitializedIterator, error) {

	logs, sub, err := _Monitor.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &MonitorInitializedIterator{contract: _Monitor.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Monitor *MonitorFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *MonitorInitialized) (event.Subscription, error) {

	logs, sub, err := _Monitor.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MonitorInitialized)
				if err := _Monitor.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Monitor *MonitorFilterer) ParseInitialized(log types.Log) (*MonitorInitialized, error) {
	event := new(MonitorInitialized)
	if err := _Monitor.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MonitorLaunchIterator is returned from FilterLaunch and is used to iterate over the raw logs and unpacked data for Launch events raised by the Monitor contract.
type MonitorLaunchIterator struct {
	Event *MonitorLaunch // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MonitorLaunchIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MonitorLaunch)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MonitorLaunch)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MonitorLaunchIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MonitorLaunchIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MonitorLaunch represents a Launch event raised by the Monitor contract.
type MonitorLaunch struct {
	Token common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterLaunch is a free log retrieval operation binding the contract event 0x35c5028003deeb5e7d5729f351ba80a2026ebbd58812af8d302ecdc4e4744f34.
//
// Solidity: event Launch(address token)
func (_Monitor *MonitorFilterer) FilterLaunch(opts *bind.FilterOpts) (*MonitorLaunchIterator, error) {

	logs, sub, err := _Monitor.contract.FilterLogs(opts, "Launch")
	if err != nil {
		return nil, err
	}
	return &MonitorLaunchIterator{contract: _Monitor.contract, event: "Launch", logs: logs, sub: sub}, nil
}

// WatchLaunch is a free log subscription operation binding the contract event 0x35c5028003deeb5e7d5729f351ba80a2026ebbd58812af8d302ecdc4e4744f34.
//
// Solidity: event Launch(address token)
func (_Monitor *MonitorFilterer) WatchLaunch(opts *bind.WatchOpts, sink chan<- *MonitorLaunch) (event.Subscription, error) {

	logs, sub, err := _Monitor.contract.WatchLogs(opts, "Launch")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MonitorLaunch)
				if err := _Monitor.contract.UnpackLog(event, "Launch", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseLaunch is a log parse operation binding the contract event 0x35c5028003deeb5e7d5729f351ba80a2026ebbd58812af8d302ecdc4e4744f34.
//
// Solidity: event Launch(address token)
func (_Monitor *MonitorFilterer) ParseLaunch(log types.Log) (*MonitorLaunch, error) {
	event := new(MonitorLaunch)
	if err := _Monitor.contract.UnpackLog(event, "Launch", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MonitorLiquityIterator is returned from FilterLiquity and is used to iterate over the raw logs and unpacked data for Liquity events raised by the Monitor contract.
type MonitorLiquityIterator struct {
	Event *MonitorLiquity // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MonitorLiquityIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MonitorLiquity)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MonitorLiquity)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MonitorLiquityIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MonitorLiquityIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MonitorLiquity represents a Liquity event raised by the Monitor contract.
type MonitorLiquity struct {
	Token    common.Address
	Reserve0 *big.Int
	Reserve1 *big.Int
	K        *big.Int
	Time     *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterLiquity is a free log retrieval operation binding the contract event 0x1955b825797bcf41aa6cbd867e8656367face6df9cc743426835a79dfcc629c0.
//
// Solidity: event Liquity(address token, uint256 reserve0, uint256 reserve1, uint256 k, uint256 time)
func (_Monitor *MonitorFilterer) FilterLiquity(opts *bind.FilterOpts) (*MonitorLiquityIterator, error) {

	logs, sub, err := _Monitor.contract.FilterLogs(opts, "Liquity")
	if err != nil {
		return nil, err
	}
	return &MonitorLiquityIterator{contract: _Monitor.contract, event: "Liquity", logs: logs, sub: sub}, nil
}

// WatchLiquity is a free log subscription operation binding the contract event 0x1955b825797bcf41aa6cbd867e8656367face6df9cc743426835a79dfcc629c0.
//
// Solidity: event Liquity(address token, uint256 reserve0, uint256 reserve1, uint256 k, uint256 time)
func (_Monitor *MonitorFilterer) WatchLiquity(opts *bind.WatchOpts, sink chan<- *MonitorLiquity) (event.Subscription, error) {

	logs, sub, err := _Monitor.contract.WatchLogs(opts, "Liquity")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MonitorLiquity)
				if err := _Monitor.contract.UnpackLog(event, "Liquity", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseLiquity is a log parse operation binding the contract event 0x1955b825797bcf41aa6cbd867e8656367face6df9cc743426835a79dfcc629c0.
//
// Solidity: event Liquity(address token, uint256 reserve0, uint256 reserve1, uint256 k, uint256 time)
func (_Monitor *MonitorFilterer) ParseLiquity(log types.Log) (*MonitorLiquity, error) {
	event := new(MonitorLiquity)
	if err := _Monitor.contract.UnpackLog(event, "Liquity", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MonitorReadyToLaunchIterator is returned from FilterReadyToLaunch and is used to iterate over the raw logs and unpacked data for ReadyToLaunch events raised by the Monitor contract.
type MonitorReadyToLaunchIterator struct {
	Event *MonitorReadyToLaunch // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MonitorReadyToLaunchIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MonitorReadyToLaunch)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MonitorReadyToLaunch)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MonitorReadyToLaunchIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MonitorReadyToLaunchIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MonitorReadyToLaunch represents a ReadyToLaunch event raised by the Monitor contract.
type MonitorReadyToLaunch struct {
	Token common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterReadyToLaunch is a free log retrieval operation binding the contract event 0x0cf95f2750965807c9d05ad42765aff0899eeb6bbe09eea63dc296c92fe671c6.
//
// Solidity: event ReadyToLaunch(address token)
func (_Monitor *MonitorFilterer) FilterReadyToLaunch(opts *bind.FilterOpts) (*MonitorReadyToLaunchIterator, error) {

	logs, sub, err := _Monitor.contract.FilterLogs(opts, "ReadyToLaunch")
	if err != nil {
		return nil, err
	}
	return &MonitorReadyToLaunchIterator{contract: _Monitor.contract, event: "ReadyToLaunch", logs: logs, sub: sub}, nil
}

// WatchReadyToLaunch is a free log subscription operation binding the contract event 0x0cf95f2750965807c9d05ad42765aff0899eeb6bbe09eea63dc296c92fe671c6.
//
// Solidity: event ReadyToLaunch(address token)
func (_Monitor *MonitorFilterer) WatchReadyToLaunch(opts *bind.WatchOpts, sink chan<- *MonitorReadyToLaunch) (event.Subscription, error) {

	logs, sub, err := _Monitor.contract.WatchLogs(opts, "ReadyToLaunch")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MonitorReadyToLaunch)
				if err := _Monitor.contract.UnpackLog(event, "ReadyToLaunch", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseReadyToLaunch is a log parse operation binding the contract event 0x0cf95f2750965807c9d05ad42765aff0899eeb6bbe09eea63dc296c92fe671c6.
//
// Solidity: event ReadyToLaunch(address token)
func (_Monitor *MonitorFilterer) ParseReadyToLaunch(log types.Log) (*MonitorReadyToLaunch, error) {
	event := new(MonitorReadyToLaunch)
	if err := _Monitor.contract.UnpackLog(event, "ReadyToLaunch", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MonitorRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the Monitor contract.
type MonitorRoleAdminChangedIterator struct {
	Event *MonitorRoleAdminChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MonitorRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MonitorRoleAdminChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MonitorRoleAdminChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MonitorRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MonitorRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MonitorRoleAdminChanged represents a RoleAdminChanged event raised by the Monitor contract.
type MonitorRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Monitor *MonitorFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*MonitorRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _Monitor.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &MonitorRoleAdminChangedIterator{contract: _Monitor.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Monitor *MonitorFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *MonitorRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _Monitor.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MonitorRoleAdminChanged)
				if err := _Monitor.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Monitor *MonitorFilterer) ParseRoleAdminChanged(log types.Log) (*MonitorRoleAdminChanged, error) {
	event := new(MonitorRoleAdminChanged)
	if err := _Monitor.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MonitorRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the Monitor contract.
type MonitorRoleGrantedIterator struct {
	Event *MonitorRoleGranted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MonitorRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MonitorRoleGranted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MonitorRoleGranted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MonitorRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MonitorRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MonitorRoleGranted represents a RoleGranted event raised by the Monitor contract.
type MonitorRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Monitor *MonitorFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*MonitorRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Monitor.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &MonitorRoleGrantedIterator{contract: _Monitor.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Monitor *MonitorFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *MonitorRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Monitor.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MonitorRoleGranted)
				if err := _Monitor.contract.UnpackLog(event, "RoleGranted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Monitor *MonitorFilterer) ParseRoleGranted(log types.Log) (*MonitorRoleGranted, error) {
	event := new(MonitorRoleGranted)
	if err := _Monitor.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MonitorRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the Monitor contract.
type MonitorRoleRevokedIterator struct {
	Event *MonitorRoleRevoked // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MonitorRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MonitorRoleRevoked)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MonitorRoleRevoked)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MonitorRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MonitorRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MonitorRoleRevoked represents a RoleRevoked event raised by the Monitor contract.
type MonitorRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Monitor *MonitorFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*MonitorRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Monitor.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &MonitorRoleRevokedIterator{contract: _Monitor.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Monitor *MonitorFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *MonitorRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Monitor.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MonitorRoleRevoked)
				if err := _Monitor.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Monitor *MonitorFilterer) ParseRoleRevoked(log types.Log) (*MonitorRoleRevoked, error) {
	event := new(MonitorRoleRevoked)
	if err := _Monitor.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MonitorSellIterator is returned from FilterSell and is used to iterate over the raw logs and unpacked data for Sell events raised by the Monitor contract.
type MonitorSellIterator struct {
	Event *MonitorSell // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MonitorSellIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MonitorSell)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MonitorSell)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MonitorSellIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MonitorSellIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MonitorSell represents a Sell event raised by the Monitor contract.
type MonitorSell struct {
	Sender     common.Address
	Token      common.Address
	AmountOIn  *big.Int
	Amount1In  *big.Int
	Amount0Out *big.Int
	Amount1Out *big.Int
	To         common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterSell is a free log retrieval operation binding the contract event 0xce0f086a8f92779a80758b1d481f2437e363e3c36bb300b37baa5b451b563dfd.
//
// Solidity: event Sell(address sender, address token, uint256 amountOIn, uint256 amount1In, uint256 amount0Out, uint256 amount1Out, address to)
func (_Monitor *MonitorFilterer) FilterSell(opts *bind.FilterOpts) (*MonitorSellIterator, error) {

	logs, sub, err := _Monitor.contract.FilterLogs(opts, "Sell")
	if err != nil {
		return nil, err
	}
	return &MonitorSellIterator{contract: _Monitor.contract, event: "Sell", logs: logs, sub: sub}, nil
}

// WatchSell is a free log subscription operation binding the contract event 0xce0f086a8f92779a80758b1d481f2437e363e3c36bb300b37baa5b451b563dfd.
//
// Solidity: event Sell(address sender, address token, uint256 amountOIn, uint256 amount1In, uint256 amount0Out, uint256 amount1Out, address to)
func (_Monitor *MonitorFilterer) WatchSell(opts *bind.WatchOpts, sink chan<- *MonitorSell) (event.Subscription, error) {

	logs, sub, err := _Monitor.contract.WatchLogs(opts, "Sell")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MonitorSell)
				if err := _Monitor.contract.UnpackLog(event, "Sell", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSell is a log parse operation binding the contract event 0xce0f086a8f92779a80758b1d481f2437e363e3c36bb300b37baa5b451b563dfd.
//
// Solidity: event Sell(address sender, address token, uint256 amountOIn, uint256 amount1In, uint256 amount0Out, uint256 amount1Out, address to)
func (_Monitor *MonitorFilterer) ParseSell(log types.Log) (*MonitorSell, error) {
	event := new(MonitorSell)
	if err := _Monitor.contract.UnpackLog(event, "Sell", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MonitorUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the Monitor contract.
type MonitorUpgradedIterator struct {
	Event *MonitorUpgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MonitorUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MonitorUpgraded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MonitorUpgraded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MonitorUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MonitorUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MonitorUpgraded represents a Upgraded event raised by the Monitor contract.
type MonitorUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Monitor *MonitorFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*MonitorUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Monitor.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &MonitorUpgradedIterator{contract: _Monitor.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Monitor *MonitorFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *MonitorUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Monitor.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MonitorUpgraded)
				if err := _Monitor.contract.UnpackLog(event, "Upgraded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Monitor *MonitorFilterer) ParseUpgraded(log types.Log) (*MonitorUpgraded, error) {
	event := new(MonitorUpgraded)
	if err := _Monitor.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
