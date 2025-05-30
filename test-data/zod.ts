import { z } from "zod";

export const UserSchema = z.object({
  id: z.number(),
  username: z.string(),
  email: z.string().email(),
  full_name: z.string().nullable(),
});

export type User = z.infer<typeof UserSchema>;

export const PostSchema = z.object({
  id: z.number(),
  title: z.string(),
  content: z.string().nullable(),
  author_id: z.number(),
});
export type Post = z.infer<typeof PostSchema>;

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
