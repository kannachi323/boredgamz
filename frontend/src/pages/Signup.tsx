import { useState } from 'react';

export default function SignUp() {
  const [formData, setFormData] = useState({
    email: '',
    password: '',
  });

    function handleChange(e : React.ChangeEvent<HTMLInputElement>) {
        const { name, value } = e.target;
        setFormData((prev) => ({ ...prev, [name]: value }));
    };

    async function handleSubmit(e : React.FormEvent<HTMLFormElement>) {
        e.preventDefault();


        if (!formData.email || !formData.password) {
        alert("Please fill out both fields");
        return;
        }

        console.log("User registered:", formData);

        const res = await fetch("http://localhost:3000/signup", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ email: formData.email, password: formData.password }),
        })

        const data = await res.json();

        console.log(data);
        
        setFormData({ email: '', password: '' });
    };

  return (
    <form onSubmit={(e) => handleSubmit(e)} style={{ maxWidth: 400, margin: '0 auto' }}>
      <h2>Sign Up</h2>

      <label style={{ display: 'block', marginBottom: 10 }}>
        Email:
        <input
          type="email"
          name="email"
          value={formData.email}
          onChange={handleChange}
          required
          style={{ display: 'block', width: '100%', padding: 8 }}
        />
      </label>

      <label style={{ display: 'block', marginBottom: 10 }}>
        Password:
        <input
          type="password"
          name="password"
          value={formData.password}
          onChange={handleChange}
          required
          style={{ display: 'block', width: '100%', padding: 8 }}
        />
      </label>

      <button type="submit" style={{ padding: '10px 20px' }}>
        Register
      </button>
    </form>
  );
}
