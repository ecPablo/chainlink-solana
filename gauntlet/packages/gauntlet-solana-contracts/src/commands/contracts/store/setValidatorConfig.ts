import { Result } from '@chainlink/gauntlet-core'
import { logger, BN } from '@chainlink/gauntlet-core/dist/utils'
import { SolanaCommand, TransactionResponse } from '@chainlink/gauntlet-solana'
import { PublicKey } from '@solana/web3.js'
import { CONTRACT_LIST, getContract } from '../../../lib/contracts'
import { getRDD } from '../../../lib/rdd'

type Input = {
  store: string
  threshold: number | string
}

export default class SetValidatorConfig extends SolanaCommand {
  static id = 'store:set_validator_config'
  static category = CONTRACT_LIST.OCR_2

  static examples = [
    'yarn gauntlet store:set_validator_config --network=devnet --state=EPRYwrb1Dwi8VT5SutS4vYNdF8HqvE7QwvqeCCwHdVLC --threshold=1000 --feed=EPRYwrb1Dwi8VT5SutS4vYNdF8HqvE7QwvqeCCwHdVLC',
  ]

  makeInput = (userInput): Input => {
    if (userInput) return userInput as Input
    const rdd = getRDD(this.flags.rdd)
    const aggregator = rdd.contracts[this.flags.state]
    return {
      store: aggregator.store,
      threshold: this.flags.threshold,
    }
  }

  constructor(flags, args) {
    super(flags, args)

    this.require(!!this.flags.state, 'Please provide flags with "state"')
  }

  execute = async () => {
    const store = getContract(CONTRACT_LIST.STORE, '')
    const address = store.programId.toString()
    const program = this.loadProgram(store.idl, address)

    const state = new PublicKey(this.flags.state)
    const input = this.makeInput(this.flags.input)
    const owner = this.wallet.payer

    console.log('INPUT', input)

    const store = new PublicKey(input.store)
    const threshold = new BN(input.threshold)

    console.log(`Setting store config on ${state.toString()}...`)

    const tx = await program.rpc.setValidatorConfig(threshold, {
      accounts: {
        state: store.publicKey,
        store: feed.publicKey,
        authority: owner.publicKey,
      },
      signers: [owner],
    })

    logger.success(`Validator config on tx ${tx}`)

    return {
      responses: [
        {
          tx: this.wrapResponse(tx, state.toString(), { state: state.toString() }),
          contract: state.toString(),
        },
      ],
    } as Result<TransactionResponse>
  }
}
