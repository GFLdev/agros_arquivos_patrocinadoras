export function validateUsername(username: string | null | undefined) {
  return !!username && username.length >= 4 && username.length <= 16
}

export function validatePassword(password: string | null | undefined) {
  return !!password && password.length >= 4
}
