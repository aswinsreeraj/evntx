function RegisterForm({ email }: any) {

  return (
    <div>

      <h2 className="text-xl font-semibold mb-4">
        Welcome to EVNTX family
      </h2>

      <input
        value={email}
        disabled
        className="w-full border rounded-lg p-3 mb-4"
      />

      <div className="grid grid-cols-2 gap-3 mb-4">
        <input placeholder="First Name" className="border p-3 rounded-lg"/>
        <input placeholder="Last Name" className="border p-3 rounded-lg"/>
      </div>

      <input
        placeholder="Date of birth"
        className="w-full border rounded-lg p-3 mb-4"
      />

      <button className="w-full bg-black text-white py-3 rounded-lg">
        Register
      </button>

    </div>
  )
}

export default RegisterForm