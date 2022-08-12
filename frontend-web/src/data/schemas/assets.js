import { schema } from 'normalizr'
import { userSchema } from './system'

const bankAccountBalanceSchema = new schema.Entity(
  'bankAccountBalances',
  {},
  {
    idAtribute: (bankAccountBalance) => bankAccountBalance.id.toLowerCase(),
  }
)

const bankAccountSchema = new schema.Entity(
  'bankAccounts',
  {
    user: userSchema,
  },
  {
    idAttribute: (bank) => bank.id.toLowerCase(),
  }
)

bankAccountBalanceSchema.define({ bankAccount: bankAccountSchema })
bankAccountSchema.define({ balances: [bankAccountBalanceSchema] })

export { bankAccountSchema, bankAccountBalanceSchema }
