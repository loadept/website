import { Code } from "lucide-preact"
import type { Post, Resources } from '../types/home'
import { ResourcesTree } from './ResourcesTree'

interface Props {
  categories: string[]
  posts: Post[]
  resources: Resources[]
}

export const Content = ({ categories, posts, resources }: Props) => {
  return (
    <section className="space-y-6 mt-15 mb-16">        
      <h2 className="flex items-center text-2xl font-bold mb-6 border-b border-border pb-3">
        <Code class="w-5 h-5 mr-2" />
        Contenido
      </h2>
      <ResourcesTree categories={categories} posts={posts} resources={resources}  />
    </section>
  )
}
