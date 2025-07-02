import { z } from "zod";

// [agree:user:zod]
export const UserSchema = z.object({
  id: z.number(),
  username: z.string(),
  email: z.string().email(),
  full_name: z.string().nullable(),
});
export type User = z.infer<typeof UserSchema>;
// [agree:end]

// [agree:post:zod]
export const PostSchema = z.object({
  id: z.number(),
  user: z.string(),
});
export type Post = z.infer<typeof PostSchema>;
// [agree:end]

// Additional schemas for extended testing
export const AddressSchema = z.object({
  id: z.number(),
  user: UserSchema,
  street: z.string(),
  city: z.string(),
  state: z.string(),
  zip_code: z.string(),
});
export type Address = z.infer<typeof AddressSchema>;

export const OrganizationSchema = z.object({
  id: z.number(),
  name: z.string(),
  domain: z.string(),
  description: z.string().nullable(),
  owner: UserSchema,
  departments: z.array(z.string()),
});
export type Organization = z.infer<typeof OrganizationSchema>;
