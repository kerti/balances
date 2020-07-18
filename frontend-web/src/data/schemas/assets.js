import { schema } from "normalizr";
import { userSchema } from "./system";

const bankAccountSchema = new schema.Entity(
  "bankAccounts",
  {
    user: userSchema,
  },
  {
    idAttribute: (bank) => bank.id.toLowerCase(),
  }
);

const bankAccountBalanceSchema = new schema.Entity(
  "bankAccountBalance",
  {
    bankAccount: bankAccountSchema,
  },
  {
    idAtribute: (bankAccountBalance) => bankAccountBalance.id.toLowerCase(),
  }
);

export { bankAccountSchema, bankAccountBalanceSchema };
