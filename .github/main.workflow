workflow "Golang workflow -> Lint" {
  on = "push"
  resolves = ["GolangCI-Lint"]
}

action "GolangCI-Lint" {
  uses = "./.github/actions/golang"
  args = "lint"
}
