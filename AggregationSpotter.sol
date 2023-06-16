//SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.19;

import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/AccessControlUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/utils/ContextUpgradeable.sol";

contract AggregationSpotter is Initializable, UUPSUpgradeable, AccessControlUpgradeable, OwnableUpgradeable {
    error AggregationSpotter__KeeperIsNotAllowed(address);
    error AggregationSpotter__ContractIsNotAllowed(address);
    error AggregationSpotter__OperationIsAlreadyApproved(uint256);
    error AggregationSpotter__KeeperIsAlreadyApproved(address, uint256);
    error AggregationSpotter__OperationDoesNotExist(uint256);
    error AggregationSpotter__CallFailed(uint256);
    error AggregationSpotter__OpIsNotApprovedOrExecuted(uint256);

    modifier isAllowedKeeper(address toCheck) {
        if (!allowedKeepers[toCheck]) {
            revert AggregationSpotter__KeeperIsNotAllowed(toCheck);
        }
        _;
    }
    modifier isAllowedContract(address toCheck) {
        if (!allowedContracts[toCheck]) {
            revert AggregationSpotter__ContractIsNotAllowed(toCheck);
        }
        _;
    }

    struct KeeperProof {
        address keeper;
        uint256 gasPrice;
    }
    /// @notice struct for informations that hold knowledge of operation status
    /// @param isApproved indicates operation approved and ready to execute
    /// @param isExecuted indicates operation is executed
    /// @param proofsCount number of proofs by uniq keepers
    /// @param proofedKeepers keepers which made proof of operation
    struct ProofInfo {
        bool isApproved;
        bool isExecuted;
        uint32 proofsCount;
        KeeperProof[] proofedKeepers;
    }

    /// @notice struct for informations that hold knowledge of operation calling proccess
    /// @param contr contract address
    /// @param functionSelector function selector to execute
    /// @param params parameters for func call
    struct OperationData {
        address contr;
        bytes4 functionSelector;
        bytes params;
    }

    /// @notice main struct that keeps all information about one entity
    /// @param proofInfo struct for informations that hold knowledge of operation status
    /// @param oracleOpTxId tx id where operation was generated on entangle oracle blockchain spotter contract
    /// @param operationData struct for informations that hold knowledge of operation calling proccess
    struct Operation {
        ProofInfo proofInfo;
        bytes32 oracleOpTxId;
        OperationData operationData;
    }

    /// @notice map with operations, key: opHash is a uint256(keccak256(uint256 oracleOpTxId, OperationData operationData))), value - Operation which need to validate and execute
    mapping(uint256 opHash => Operation operation) operations;

    /// @notice map of allowed contracts which we can interact
    mapping(address => bool) allowedContracts;

    /// @notice map of allowed keepers which can propose operation
    mapping(address => bool) allowedKeepers;

    uint256 numberOfAllowedKeepers;

    /// @notice 10000 = 100%
    uint256 constant rateDecimals = 10000;

    /// @notice percentage of proofs div numberOfAllowedKeepers which should be reached to approve operation. Scaled with 10000 decimals, e.g. 6000 is 60%
    uint256 consensusTargetRate;

    event NewOperation(uint256 opHash, address cont, bytes4 functionSelector);
    event ProposalApproved(uint256 opHash);
    event ProposalExecuted(uint256 opHash);

    bytes32 public constant ADMIN = keccak256("ADMIN");
    bytes32 public constant EXECUTOR = keccak256("EXECUTOR");

    function initialize(address[2] calldata initAddr) public initializer {
        _grantRole(ADMIN, initAddr[0]);
        _grantRole(EXECUTOR, initAddr[1]);
    }

    function _authorizeUpgrade(address) internal override onlyOwner {}

    /// @notice Adding keeper to whitelist
    /// @param keeper address of keeper to add
    function addKeeper(address keeper) external onlyRole(ADMIN) {
        allowedKeepers[keeper] = true;
        numberOfAllowedKeepers++;
    }

    /// @notice Removing keeper from whitelist
    /// @param keeper address of keeper to remove
    function removeKeeper(address keeper) external onlyRole(ADMIN) {
        allowedKeepers[keeper] = false;
        numberOfAllowedKeepers--;
    }

    /// @notice Adding contract to whitelist
    /// @param _contract address of contract to add
    function addAllowedContract(address _contract) external onlyRole(ADMIN) {
        allowedContracts[_contract] = true;
    }

    /// @notice Removing contract from whitelist
    /// @param _contract address of contract to remove
    function removeAllowedContract(address _contract) external onlyRole(ADMIN) {
        allowedContracts[_contract] = false;
    }

    /// @notice Setting of target rate
    /// @param rate target rate
    function setConsensusTargetRate(uint256 rate) external onlyRole(ADMIN) {
        consensusTargetRate = rate;
    }

    /// @notice proposing an operation/approve an operation/give an operation of status approved
    /// @param oracleOpTxId 1
    /// @param opData 2
    function proposeOperation(
        bytes32 oracleOpTxId,
        OperationData memory opData
    ) external isAllowedKeeper(_msgSender()) isAllowedContract(opData.contr) {
        uint256 opHash = uint256(keccak256(abi.encode(oracleOpTxId, opData)));
        address msgSender = _msgSender();
        if (operations[opHash].operationData.contr == address(0)) {
            KeeperProof[] memory kp = new KeeperProof[](1);
            kp[0] = KeeperProof(msgSender, tx.gasprice);
            Operation memory op = Operation(ProofInfo(false, false, 1, kp), oracleOpTxId, opData);
            operations[opHash] = op;
            emit NewOperation(opHash, opData.contr, opData.functionSelector);
        } else {
            if (operations[opHash].proofInfo.isApproved) {
                revert AggregationSpotter__OperationIsAlreadyApproved(opHash);
            }
            for (uint256 i = 0; i < operations[opHash].proofInfo.proofedKeepers.length; i++) {
                if (operations[opHash].proofInfo.proofedKeepers[i].keeper == msgSender) {
                    revert AggregationSpotter__KeeperIsAlreadyApproved(msgSender, opHash);
                }
            }
            operations[opHash].proofInfo.proofedKeepers.push(KeeperProof(msgSender, tx.gasprice));
            operations[opHash].proofInfo.proofsCount++;
            if (
                (operations[opHash].proofInfo.proofsCount * rateDecimals) / numberOfAllowedKeepers >=
                consensusTargetRate
            ) {
                operations[opHash].proofInfo.isApproved = true;
                emit ProposalApproved(opHash);
            }
        }
    }

    /// @notice execute approved operation
    /// @param opHash 1
    function executeOperation(uint256 opHash) external onlyRole(EXECUTOR) {
        if (operations[opHash].operationData.contr == address(0)) {
            revert AggregationSpotter__OperationDoesNotExist(opHash);
        }
        if (operations[opHash].proofInfo.isApproved && !operations[opHash].proofInfo.isExecuted) {
            (bool success, ) = address(operations[opHash].operationData.contr).call(
                abi.encodeWithSelector(
                    operations[opHash].operationData.functionSelector,
                    operations[opHash].operationData.params
                )
            );
            if (!success) {
                revert AggregationSpotter__CallFailed(opHash);
            }
            operations[opHash].proofInfo.isExecuted = true;
            emit ProposalExecuted(opHash);
        } else {
            revert AggregationSpotter__OpIsNotApprovedOrExecuted(opHash);
        }
    }
}
