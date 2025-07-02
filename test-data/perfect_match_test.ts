// Test schemas that should be perfectly equivalent cross-language

// [agree:perfect:zod]
export const PerfectSchema = z.object({
  id: z.number(),        // Should match int/integer
  name: z.string(),      // Should match str/string
  active: z.boolean(),   // Should match bool/boolean
  score: z.number(),     // Should match float (both map to number)
  email: z.string().email(), // Should match EmailStr (both map to email)
});
export type Perfect = z.infer<typeof PerfectSchema>;
// [agree:end]