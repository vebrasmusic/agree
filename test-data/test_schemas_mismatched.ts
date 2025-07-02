// Test data with intentionally mismatched schemas for testing

// [agree:mismatch_test:zod]
export const MismatchTestSchema = z.object({
  id: z.number(),
  name: z.string(),
  email: z.string().email(), // Different email validation
  score: z.number(),  // This field missing in Python
  // missing 'age' field that exists in Python
});
export type MismatchTest = z.infer<typeof MismatchTestSchema>;
// [agree:end]