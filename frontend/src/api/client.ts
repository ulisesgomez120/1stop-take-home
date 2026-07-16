const baseURL = import.meta.env.VITE_API_BASE_URL

export async function apiGet<T>(path: string): Promise<T> {
  const res = await fetch(`${baseURL}${path}`)
  if (!res.ok) {
    throw new Error(`GET ${path} failed: ${res.status}`)
  }
  return res.json() as Promise<T>
}

export async function apiPut<T>(path: string, body: unknown): Promise<T> {
  const res = await fetch(`${baseURL}${path}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })
  if (!res.ok) {
    throw new Error(`PUT ${path} failed: ${res.status}`)
  }
  return res.json() as Promise<T>
}
