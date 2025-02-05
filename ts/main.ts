import Neon from "@cityofzion/neon-js"
import { Client, Token } from "./itdl"

const fromAccount = Neon.create.account(
	"L557bcBh8mXTp68ABb8Bt4VzyMPAgQeyXib3s2KRot1TCproLGDV"
)
const toAddress = "NS7JVaxBCuS4ffYKWaWjbYS9jTKQaaJxnB"

const token = new Token(fromAccount)
const hash = await token.transfer(toAddress, 10000)

console.log(hash)

const client = new Client()

console.log(await client.balance(fromAccount.address))
console.log(await client.balance(toAddress))

const txs = await client.transactions(fromAccount.address)

for (const tx of txs)
	console.log(
		tx.type === "sent" ? "-->" : "<--",
		tx.amount,
		tx.address || "<nil>",
		tx.blockindex
	)
