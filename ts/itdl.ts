import type { wallet } from "@cityofzion/neon-core"
import type { Nep17TransferEvent } from "@cityofzion/neon-core/lib/rpc"
import Neon, { u, sc } from "@cityofzion/neon-js"
import type { SmartContract } from "@cityofzion/neon-js/lib/experimental"

const gas = Neon.CONST.NATIVE_CONTRACT_HASH.GasToken
const hexgas = `0x${gas}`
const mainToken = u.HexString.fromHex(gas)

const networkMagic = 56753
const rpcAddress = "http://localhost:20331"

type Tx = {
	type: "received" | "sent"
	amount: string
	blockindex: number
	timestamp: number
	address?: string
	hash: string
}

export class Client {
	#client = Neon.create.rpcClient(rpcAddress)

	async balance(address: string) {
		const balances = await this.#client.getNep17Balances(address)

		for (const balance of balances.balance)
			if (balance.assethash === hexgas) return balance.amount

		return 0
	}

	async transactions(address: string, type?: Tx["type"]) {
		const txs = await this.#client.getNep17Transfers(address, "0")

		const result: Tx[] = []

		const txm = (t: Tx["type"], tx: Nep17TransferEvent): Tx => ({
			type: t,
			amount: tx.amount,
			blockindex: tx.blockindex,
			timestamp: tx.timestamp,
			address: tx.transferaddress,
			hash: tx.txhash,
		})

		if (!type || type === "received")
			for (const tx of txs.received)
				if (tx.assethash === hexgas) result.push(txm("received", tx))

		if (!type || type === "sent")
			for (const tx of txs.sent)
				if (tx.assethash === hexgas) result.push(txm("sent", tx))

		return result
	}
}

export class Token {
	account: wallet.Account
	#contract: SmartContract

	constructor(account: wallet.Account) {
		this.account = account
		this.#contract = new Neon.experimental.SmartContract(mainToken, {
			networkMagic,
			rpcAddress,
			account,
		})
	}

	transfer(toAddress: string, amount: number) {
		const params = [
			sc.ContractParam.hash160(this.account.address),
			sc.ContractParam.hash160(toAddress),
			sc.ContractParam.integer(amount),
			sc.ContractParam.any(null),
		]

		return this.#contract.invoke("transfer", params)
	}
}
