import { defineCollection, z } from 'astro:content'
import { glob } from 'astro/loaders'

const posts = defineCollection({
  loader: glob({ pattern: 'posts/**/*.md', base: './src/content' }),
  schema: z.object({
    title: z.string(),
    date: z.coerce.date(),
    keywords: z.array(z.string()),
    category: z.string().optional()
  })
})

export const collections = {
  posts
}
