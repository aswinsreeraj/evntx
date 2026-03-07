import { useState } from "react"
import { motion, AnimatePresence } from "framer-motion"

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
    <Modal open={open} onClose={onClose} className="bg-white rounded-xl w-[900px] h-[550px] relative overflow-hidden flex">
      {/* Container holding image and content */}
      <div className="flex w-full h-full">

        {/* Left image - Full height */}
        <div className="hidden md:block w-1/2 h-full bg-gray-100">
          <img
            src="https://images.unsplash.com/photo-1543269865-cbf427effbad?w=800&q=80"
            alt="Group of friends"
            className="h-full w-full object-cover"
          />
        </div>

        {/* Animated content */}
        <div className="relative w-full md:w-1/2 p-10 flex flex-col justify-center overflow-x-hidden">

          <AnimatePresence mode="wait">

            {view === "login-choice" && (
              <motion.div
                key="login-choice"
                initial={{ x: 50, opacity: 0 }}
                animate={{ x: 0, opacity: 1 }}
                exit={{ x: -50, opacity: 0 }}
                transition={{ duration: 0.25 }}
              >
                <LoginChoice
                  setView={setView}
                  setEmail={setEmail}
                />
              </motion.div>
            )}

            {view === "email-input" && (
              <motion.div
                key="email-input"
                initial={{ x: 50, opacity: 0 }}
                animate={{ x: 0, opacity: 1 }}
                exit={{ x: -50, opacity: 0 }}
                transition={{ duration: 0.25 }}
              >
                <EmailInput
                  email={email}
                  setEmail={setEmail}
                  setView={setView}
                />
              </motion.div>
            )}

            {view === "otp-verify" && (
              <motion.div
                key="otp-verify"
                initial={{ x: 50, opacity: 0 }}
                animate={{ x: 0, opacity: 1 }}
                exit={{ x: -50, opacity: 0 }}
                transition={{ duration: 0.25 }}
              >
                <OTPVerify
                  email={email}
                  setView={setView}
                />
              </motion.div>
            )}

            {view === "register" && (
              <motion.div
                key="register"
                initial={{ x: 50, opacity: 0 }}
                animate={{ x: 0, opacity: 1 }}
                exit={{ x: -50, opacity: 0 }}
                transition={{ duration: 0.25 }}
              >
                <RegisterForm email={email} />
              </motion.div>
            )}

          </AnimatePresence>

        </div>
      </div>

      {view !== "login-choice" && (
        <button
          onClick={() => setView("login-choice")}
          className="absolute top-6 left-6 text-gray-400 hover:text-gray-600 transition-colors"
        >
          ← Back
        </button>
      )}
    </Modal >
  )
}

export default AuthModal