import { schema } from 'normalizr'
import { userSchema } from './system'

const bankAccountBalanceSchema = new schema.Entity(
  'bankAccountBalances',
  {
    bankAccount: bankAccountSchema,
  },
  {
    idAtribute: (bankAccountBalance) => bankAccountBalance.id.toLowerCase(),
  }
)

const bankAccountSchema = new schema.Entity(
  'bankAccounts',
  {
    user: userSchema,
    balances: [bankAccountBalanceSchema],
  },
  {
    idAttribute: (bank) => bank.id.toLowerCase(),
  }
)

export { bankAccountSchema, bankAccountBalanceSchema }
