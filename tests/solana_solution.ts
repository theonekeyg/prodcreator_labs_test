import * as anchor from "@coral-xyz/anchor";
import { Program } from "@coral-xyz/anchor";
import { expect } from "chai";
import { AggregationSpotter, TestLinker } from "../target/types/solana_solution";

const KEEPER_SEED = "keeper";
const CONTRACT_SEED = "contract";
const OPERATION_SEED = "operation";

async function get_keeper_pda(keeper: PublicKey, programId: PublicKey) {
  const [keeper_rv, bump] = await anchor.web3.PublicKey.findProgramAddress([
    anchor.utils.bytes.utf8.encode(KEEPER_SEED),
    keeper.toBytes()
  ],
  programId);
  return [keeper_rv, bump]
}

async function get_contract_pda(contract: PublicKey, programId: PublicKey) {
  const [contract_rv, bump] = await anchor.web3.PublicKey.findProgramAddress([
    anchor.utils.bytes.utf8.encode(CONTRACT_SEED),
    contract.toBytes()
  ],
  programId);
  return [contract_rv, bump]
}

async function get_operation_pda(operation: PublicKey, programId: PublicKey) {
  const [operation_rv, bump] = await anchor.web3.PublicKey.findProgramAddress([
    anchor.utils.bytes.utf8.encode(OPERATION_SEED),
    operation.toBytes()
  ],
  programId);
  return [operation_rv, bump]
}

interface OperationData {
  contr: PublicKey,
}

describe("solana_solution", async () => {
  // Configure the client to use the local cluster.
  let provider_opts = {
    // preflightCommitment: "processed",
    // commitment: "processed",
    logLevel: "debug",
  };
  anchor.setProvider(anchor.AnchorProvider.env(provider_opts));

  const program = anchor.workspace.AggregationSpotter as Program<AggregationSpotter>;
  const testLinker = anchor.workspace.TestLinker as Program<TestLinker>;

  const admin = program.provider.wallet.payer;
  const executor = anchor.web3.Keypair.generate();
  const spotter = anchor.web3.Keypair.generate();

  const keeper1_acc = anchor.web3.Keypair.generate();
  const keeper2_acc = anchor.web3.Keypair.generate();
  const keeper3_acc = anchor.web3.Keypair.generate();
  const [keeper1, ] = await get_keeper_pda(keeper1_acc.publicKey, program.programId);
  const [keeper2, ] = await get_keeper_pda(keeper2_acc.publicKey, program.programId);
  const [keeper3, ] = await get_keeper_pda(keeper3_acc.publicKey, program.programId);

  const contract_acc = anchor.web3.Keypair.generate();
  const [contract, ] = await get_contract_pda(contract_acc.publicKey, program.programId);

  const operation_acc = anchor.web3.Keypair.generate();
  const [operation, ] = await get_operation_pda(operation_acc.publicKey, program.programId);

  it("initializes the program", async () => {
    // Initialize the program
    // Don't pass admin signature, since it's already the payer. This is the case for all tests
    await program.methods
      .initialize([admin.publicKey, executor.publicKey])
      .accounts({
        spotter: spotter.publicKey,
        admin: admin.publicKey,
      })
      .signers([spotter])
      .rpc();

    // Fetch current spotter state
    let spotterState = await program.account.aggregationSpotter.fetch(spotter.publicKey);
    expect(spotterState.isInitialized).to.eql(true);
    expect(spotterState.admin).to.eql(admin.publicKey);
    expect(spotterState.executor).to.eql(executor.publicKey);

  });

  it("add keeper", async () => {
    await program.methods
      .createKeeper(keeper1_acc.publicKey)
      .accounts({
        spotter: spotter.publicKey,
        admin: admin.publicKey,
        keeperPda: keeper1,
        keeperAcc: keeper1_acc.publicKey,
      })
      .signers([])
      .rpc();

    let keeperState = await program.account.keeper.fetch(keeper1);
    let spotterState = await program.account.aggregationSpotter.fetch(spotter.publicKey);
    expect(keeperState.isAllowed).to.eql(true);
    expect(keeperState.key).to.eql(keeper1_acc.publicKey);
    // expect(spotterState.numberOfKeepers).to.eql(new anchor.BN(1));
  });

  it("remove keeper and add 3 other keepers", async () => {
    await program.methods
      .removeKeeper(keeper1_acc.publicKey)
      .accounts({
        spotter: spotter.publicKey,
        admin: admin.publicKey,
        keeperPda: keeper1,
        keeperAcc: keeper1_acc.publicKey,
      })
      .signers([])
      .rpc();

    let keeperState = await program.account.keeper.fetch(keeper1);
    expect(keeperState.isAllowed).to.eql(false);

    await program.methods
      .enableKeeper(keeper1_acc.publicKey)
      .accounts({
        spotter: spotter.publicKey,
        admin: admin.publicKey,
        keeperPda: keeper1,
        keeperAcc: keeper1_acc.publicKey,
      })
      .signers([])
      .rpc();

    keeperState = await program.account.keeper.fetch(keeper1);
    expect(keeperState.isAllowed).to.eql(true);

    let keeper_acc: anchor.web3.Keypair;
    let keeper: PublicKey;

    for ([keeper_acc, keeper] of [[keeper2_acc, keeper2], [keeper3_acc, keeper3]]) {
      await program.methods
      .createKeeper(keeper_acc.publicKey)
      .accounts({
        spotter: spotter.publicKey,
        admin: admin.publicKey,
        keeperPda: keeper,
        keeperAcc: keeper_acc.publicKey,
      })
      .signers([])
      .rpc();
    }
  });

  it("add new allowed contract", async () => {
    await program.methods
      .addAllowedContract(contract_acc.publicKey)
      .accounts({
        spotter: spotter.publicKey,
        admin: admin.publicKey,
        contractPda: contract,
        contractAcc: contract_acc.publicKey,
      })
      .signers([])
      .rpc();

    let contractState = await program.account.allowedContract.fetch(contract);
    expect(contractState.isAllowed).to.eql(true);
  });

  it("create new operation", async () => {
    let operationData: OperationData = {
      contr: testLinker.programId,
      contracts: [contract_acc.publicKey, executor.publicKey],
    };
    let gas_price = new anchor.BN(50000);

    // Transfer some sol to the keeper account
    await program.provider.connection.confirmTransaction(
      await program.provider.connection.requestAirdrop(keeper1_acc.publicKey, 10000000000),
      "confirmed"
    );

    // Create new operation, fetch current state and check its fields
    await program.methods
      .createOperation(operationData, gas_price)
      .accounts({
        spotter: spotter.publicKey,
        keeperAcc: keeper1_acc.publicKey,
        keeperPda: keeper1,
        operationAcc: operation_acc.publicKey,
        operationPda: operation,
      })
      .signers([keeper1_acc])
      .rpc();

    let operationState = await program.account.operation.fetch(operation);

    // expect(operationState.proofInfo.proofsCount).to.eql(new anchor.BN(1));
    expect(operationState.proofInfo.isApproved).to.eql(false);
    expect(operationState.proofInfo.isExecuted).to.eql(false);
    expect(operationState.operationData.contr).to.eql(testLinker.programId);

  });

  it("propose proof on other keepers", async () => {

    let gas_price = new anchor.BN(50000);
    await program.methods
      .proposeOperation(gas_price)
      .accounts({
        spotter: spotter.publicKey,
        keeperAcc: keeper2_acc.publicKey,
        keeperPda: keeper2,
        operationAcc: operation_acc.publicKey,
        operationPda: operation,
      })
      .signers([keeper2_acc])
      .rpc();

    let operationState = await program.account.operation.fetch(operation);
    expect(operationState.proofInfo.isApproved).to.eql(true);
  });

  it("execute approved operation", async () => {

    await program.methods
      .executeOperation()
      .accounts({
        testProgram: testLinker.programId,
        executer: executor.publicKey,
        operationAcc: operation_acc.publicKey,
        operationPda: operation,
      })
      .signers([executor])
      .rpc();

    let operationState = await program.account.operation.fetch(operation);
    expect(operationState.proofInfo.isExecuted).to.eql(true);

  });
});
