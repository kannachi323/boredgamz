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
