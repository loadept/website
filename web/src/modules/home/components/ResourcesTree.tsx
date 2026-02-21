import { useState } from 'preact/hooks'
import type { Post, Resources } from '../types/home'

export const ResourcesTree = (
  { categories, posts, resources }:
  { categories: string[], posts: Post[], resources: Resources[] }
) => {
  const [expandedArticleCategories, setExpandedArticleCategories] = useState<string[]>([])

  const toggleArticleCategory = (categoryId: string) => {
    setExpandedArticleCategories((prev) =>
      prev.includes(categoryId) ? prev.filter((c) => c !== categoryId) : [...prev, categoryId],
    )
  }

  return (
    <div className="space-y-4">
      {resources.map((resource) => (
        <div key={resource.label}>
          <h2 className="text-primary font-normal text-lg">
            {resource.label}
          </h2>

          <ul className="ml-6 mt-2 space-y-1">
            {resource.label === "Notas"
              ?
                categories.map((category, k) => {
                  const filteredPosts = posts.filter((a) => a.category === category)
                  return (
                    <li key={k}>
                      <button
                        onClick={() => toggleArticleCategory(category)}
                        className="text-foreground hover:text-primary transition-colors font-normal cursor-pointer"
                      >
                        {expandedArticleCategories.includes(category) ? "▼" : "▶"} {category}
                      </button>

                      {expandedArticleCategories.includes(category) && (
                        <ul className="ml-6 mt-2 space-y-1">
                          {filteredPosts.map((post, k) => (
                            <li key={k}>
                              <a href={`/posts/${post.slug}`} className="hover:underline">
                                {post.title}
                              </a>
                            </li>
                          ))}
                        </ul>
                      )}
                    </li>
                  )
                })
              : 
                resource.items.map((item) => (
                  <li>
                    <a
                      key={item.label}
                      href={item.url || "#"}
                      target="_blank"
                      className="flex text-foreground font-normal hover:underline"
                    >
                      {item.label}
                    </a>
                      <p className="text-muted-foreground text-sm">{item.description}</p>
                  </li>
                ))
            }
          </ul>
        </div>
      ))}
    </div>
  )
}
