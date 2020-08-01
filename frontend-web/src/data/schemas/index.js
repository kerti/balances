import { userSchema } from './system'
import { bankAccountSchema, bankAccountBalanceSchema } from './assets'

// Schemas for Balances API responses
export const Schemas = {
  // Assets - Bank Accounts
  BANK_ACCOUNT: bankAccountSchema,
  BANK_ACCOUNT_ARRAY: [bankAccountSchema],
  BANK_ACCOUNT_BALANCE: bankAccountBalanceSchema,
  BANK_ACCOUNT_BALANCE_ARRAY: [bankAccountBalanceSchema],
  // System
  USER: userSchema,
  USER_ARRAY: [userSchema],
}
