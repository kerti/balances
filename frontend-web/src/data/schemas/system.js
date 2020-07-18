import { schema } from "normalizr";

const userSchema = new schema.Entity(
  "users",
  {},
  {
    idAttribute: (user) => user.login.toLowerCase(),
  }
);

export { userSchema };
