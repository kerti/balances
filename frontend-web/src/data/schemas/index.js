import { userSchema } from './system'
import { bankAccountSchema, bankAccountBalanceSchema } from './assets'

// Schemas for Balances API responses
export const Schemas = {
  // System
  USER: userSchema,
  USER_ARRAY: [userSchema],
}
