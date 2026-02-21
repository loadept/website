import type { FrontMatter } from "../../shared/types/frontmatter"

export interface Post extends FrontMatter {
  slug: string
}

export interface ResourceItem {
  label: string
  url: string
  description: string
}

export interface Resources {
  label: string
  items: ResourceItem[]
}
