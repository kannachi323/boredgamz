interface Props {
    email: string,
    password: string,
}

export async function signUp({ email, password } : Props) {
    const res = await fetch("http://localhost:3000/auth/sign-up", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ email: email, password: password }),
    });

    const data = await res.json();

    console.log(data);

}

async function handleLogout(setIsAuthenticated) {
    /* Implement this later....
  const res = await fetch("http://localhost:3000/logout", {
    method: "POST",
    credentials: "include",
  });

  if (res.ok) {
    setIsAuthenticated(false);
    setUser(null);
    navigate('/login'); // or home
  } else {
    alert("Failed to log out");
  }
    */

  setIsAUthenticated(false);
  setUser(null);
}
