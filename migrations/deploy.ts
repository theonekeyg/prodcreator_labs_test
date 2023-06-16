// Migrations are an early feature. Currently, they're nothing more than this
// single deploy script that's invoked from the CLI, injecting a provider
// configured from the workspace's Anchor.toml.

const anchor = require("@coral-xyz/anchor");

module.exports = async function (provider) {
  // Configure client to use the provider.
  anchor.setProvider(provider);

  // Add your deploy script here.

  const keeper1 = "AwZYmCuxVkUBQLYPrZVYBmfChUAPgD6FZzLTENtaAoAk";
  const keeper2 = "2L4B2vWxo3AXw5ekUWnngWYqohE6MyVMzYSFzK5ZbrLj";
  const keeper3 = "9XZ5zYSifV4kXNAYHBTvPUXaK57KbH7NHzFaTvXc3cZp";
  const spotter = "b4ZqA4xp7fgmxmE7aJfJbfj5HqMZPRdiRTfz9DpJvki";

  const program = anchor.workspace.AggregationSpotter;

  const admin = program.provider.wallet.payer;
  const executor = program.provider.wallet.payer;

  // let res = await program.methods
  //   .initialize([admin.publicKey, executor.publicKey])
  //   .accounts({
  //     spotter: spotter,
  //     admin: admin.publicKey,
  //   })
  //   .signers([admin])
  //   .rpc();

  let res = await program.rpc.initialize([admin.publicKey, executor.publicKey], {
    accounts: {
      spotter: spotter,
      admin: admin.publicKey,
      systemProgram: anchor.web3.SystemProgram.programId,
    },
    signers: [admin],
  });

  console.log("res: ", res);

};
