import { useState } from "react"
import Modal from "../../../shared/ui/Modal"

import LoginChoice from "./LoginChoice"
import EmailInput from "./EmailInput"
import OTPVerify from "./OTPVerify"
import RegisterForm from "./RegisterForm"

type AuthView =
  | "login-choice"
  | "email-input"
  | "otp-verify"
  | "register"

function AuthModal({ open, onClose }: any) {
  const [view, setView] = useState<AuthView>("login-choice")
  const [email, setEmail] = useState("")

  return (
    <Modal open={open} onClose={onClose}>
      <div className="grid grid-cols-2 w-[900px]">

        {/* Left image */}
        {/* <div className="hidden md:block">
          <img
            src="/assets/images/login-modal.png"
            className="h-full w-full object-cover rounded-l-xl"
          />
        </div> */}

        {/* Right content */}
        <div className="p-8">

          {view === "login-choice" && (
            <LoginChoice
              setView={setView}
              setEmail={setEmail}
            />
          )}

          {view === "email-input" && (
            <EmailInput
              email={email}
              setEmail={setEmail}
              setView={setView}
            />
          )}

          {view === "otp-verify" && (
            <OTPVerify
              email={email}
              setView={setView}
            />
          )}

          {view === "register" && (
            <RegisterForm email={email} />
          )}

        </div>

      </div>
    </Modal>
  )
}

export default AuthModal