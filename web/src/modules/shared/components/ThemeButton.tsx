import { useEffect, useState } from "preact/hooks"
import {
  Moon,
  Sun,
} from "lucide-preact"
import { getTheme, setTheme } from "../lib/theme"

export const ThemeButton = () => {
  const [theme, setThemeState] = useState<"dark" | "light">("light")

  useEffect(() => {
    setThemeState(getTheme())

    const handler = (e: any) => setThemeState(e.detail)

    window.addEventListener("theme-change", handler)
    return () => window.removeEventListener("theme-change", handler)
  }, [])

  const toggle = () => {
    setTheme(theme === "dark" ? "light" : "dark")
  }

  return (
    <button
      onClick={toggle}
      class="text-primary hover:text-accent flex items-center gap-2 transition-colors mx-6 cursor-pointer"
      aria-label="toggle theme"
      title="toggle theme"
    >
    {theme === "dark"
      ? <Sun class="h-4 w-4" />
      : <Moon class="h-4 w-4" />
    }
    </button>
  )
}
