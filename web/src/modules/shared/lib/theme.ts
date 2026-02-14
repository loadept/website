export const getTheme = (): "dark" | "light" => {
  return localStorage.getItem("theme") === "dark"
    ? "dark"
    : "light"
}

export const setTheme = (theme: "dark" | "light") => {
  localStorage.setItem("theme", theme)

  document.documentElement.classList.toggle(
    "dark",
    theme === "dark"
  )

  window.dispatchEvent(new CustomEvent("theme-change", {
    detail: theme
  }))
}
