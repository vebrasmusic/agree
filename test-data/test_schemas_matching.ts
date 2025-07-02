// Test data with schemas that should match Python counterparts

// [agree:match_test:zod]
export const MatchTestSchema = z.object({
  id: z.number(),
  name: z.string(),
  email: z.string(),
  active: z.boolean(),
});
export type MatchTest = z.infer<typeof MatchTestSchema>;
// [agree:end]