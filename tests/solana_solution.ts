import * as anchor from "@coral-xyz/anchor";
import { Program } from "@coral-xyz/anchor";
import { AggregationSpotter } from "../target/types/solana_solution";

describe("solana_solution", () => {
  // Configure the client to use the local cluster.
  anchor.setProvider(anchor.AnchorProvider.env());

  const program = anchor.workspace.AggregationSpotter as Program<AggregationSpotter>;

  it("initializes the program", async () => {
    const admin = program.provider.wallet.payer;
    const executor = anchor.web3.Keypair.generate();
    const spotter = anchor.web3.Keypair.generate();

    // Add your test here.
    let tx = await program.methods
      .initialize([admin.publicKey, executor.publicKey])
      .accounts({
        spotter: spotter.publicKey,
        admin: admin.publicKey,
      })
      .signers([spotter])
      .rpc();

    console.log("Your transaction signature", tx);
  })
});
