export async function fetchPortfolio() {
  const res = await fetch('http://localhost:8080/api/portfolio');
  if (!res.ok) {
    throw new Error('Network response was not ok');
  }
  const data = await res.json();
  return data.projects || [];
}
