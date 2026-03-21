export type FeedbackRecord = {
  rating: number
  comment: string
}

const STORAGE_KEY = "user-booking-feedback"

export const getFeedbackMap = (): Record<string, FeedbackRecord> => {
  const raw = localStorage.getItem(STORAGE_KEY)
  if (!raw) return {}

  try {
    return JSON.parse(raw) as Record<string, FeedbackRecord>
  } catch {
    localStorage.removeItem(STORAGE_KEY)
    return {}
  }
}

export const saveFeedbackMap = (feedback: Record<string, FeedbackRecord>) => {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(feedback))
}
