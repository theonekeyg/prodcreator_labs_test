// Migrations are an early feature. Currently, they're nothing more than this
// single deploy script that's invoked from the CLI, injecting a provider
// configured from the workspace's Anchor.toml.

const anchor = require("@coral-xyz/anchor");

const SPOTTER_SEED = "spotter";
const KEEPER_SEED = "keeper";
const CONTRACT_SEED = "contract";

async function get_keeper_pda(keeper, programId) {
  const [keeper_rv, bump] = await anchor.web3.PublicKey.findProgramAddress([
    anchor.utils.bytes.utf8.encode(KEEPER_SEED),
    keeper.toBytes()
  ],
  programId);
  return [keeper_rv, bump]
}

async function get_spotter_pda(programId) {
  const [spotter_rv, bump] = await anchor.web3.PublicKey.findProgramAddress([
    anchor.utils.bytes.utf8.encode(SPOTTER_SEED),
  ],
  programId);
  return [spotter_rv, bump]
}

async function get_contract_pda(contract, programId) {
  const [contract_rv, bump] = await anchor.web3.PublicKey.findProgramAddress([
    anchor.utils.bytes.utf8.encode(CONTRACT_SEED),
    contract.toBytes()
  ],
  programId);
  return [contract_rv, bump]
}

module.exports = async function (provider) {
  // Configure client to use the provider.
  anchor.setProvider(provider);
  const program = anchor.workspace.AggregationSpotter;

  // Add your deploy script here.

  const keeper1Acc = new anchor.web3.PublicKey("AwZYmCuxVkUBQLYPrZVYBmfChUAPgD6FZzLTENtaAoAk");
  const keeper2Acc = new anchor.web3.PublicKey("2L4B2vWxo3AXw5ekUWnngWYqohE6MyVMzYSFzK5ZbrLj");
  const keeper3Acc = new anchor.web3.PublicKey("9XZ5zYSifV4kXNAYHBTvPUXaK57KbH7NHzFaTvXc3cZp");
  const allowedContractAcc = new anchor.web3.PublicKey("FBXArrxAnAZpmGmVYD8jQVauohCrXHaJBgTXzUPKRyHa");

  const [keeper1, ] = await get_keeper_pda(keeper1Acc, program.programId);
  const [keeper2, ] = await get_keeper_pda(keeper2Acc, program.programId);
  const [keeper3, ] = await get_keeper_pda(keeper3Acc, program.programId);
  const [spotter, ] = await get_spotter_pda(program.programId);
  const [allowedContract, ] = await get_contract_pda(allowedContractAcc, program.programId);

  const admin = program.provider.wallet.payer;
  const executor = program.provider.wallet.payer;

  // Initialize the program
  await program.methods
    .initialize([admin.publicKey, executor.publicKey])
    .accounts({
      spotter: spotter,
      admin: admin.publicKey,
    })
    .signers([]) // Admin is already a signer
    .rpc();

  console.log(`Initialized program at ${program.programId.toBase58()}, spotter: ${spotter.toBase58()}`);

  // Create keepers
  await program.methods
    .createKeeper(keeper1Acc)
    .accounts({
      spotter: spotter,
      admin: admin.publicKey,
      keeperPda: keeper1,
      keeperAcc: keeper1Acc,
    })
    .signers([])
    .rpc();

  console.log("Added keeper: ", keeper1Acc.toBase58());

  await program.methods
    .createKeeper(keeper2Acc)
    .accounts({
      spotter: spotter,
      admin: admin.publicKey,
      keeperPda: keeper2,
      keeperAcc: keeper2Acc,
    })
    .signers([])
    .rpc();
  console.log("Added keeper: ", keeper2Acc.toBase58());

  await program.methods
    .createKeeper(keeper3Acc)

    .accounts({
      spotter: spotter,
      admin: admin.publicKey,
      keeperPda: keeper3,
      keeperAcc: keeper3Acc,
    })
    .signers([])
    .rpc();

  console.log("Added keeper: ", keeper3Acc.toBase58());

  // Add allowed contract
  await program.methods
    .addAllowedContract(allowedContractAcc)
    .accounts({
      spotter: spotter,
      admin: admin.publicKey,
      contractPda: allowedContract,
      contractAcc: allowedContractAcc,
    })
    .signers([])
    .rpc();

  console.log("Added allowed contract: ", allowedContractAcc.toBase58());

  // Fund keepers some native SOL through solana transfer from current payer
  // Transfer 5 SOL to each keeper
  for (let keeper of [keeper1Acc, keeper2Acc, keeper3Acc]) {

    const lamports = 5 * 1000000000;

    const tx = new anchor.web3.Transaction().add(
      anchor.web3.SystemProgram.transfer({
        fromPubkey: admin.publicKey,
        toPubkey: keeper,
        lamports,
      })
    );
    await provider.sendAndConfirm(tx, [admin]);

    console.log(`[tx: ${tx}] Funded 5 SOL to keeper: ${keeper.toBase58()}`);
  }
};
