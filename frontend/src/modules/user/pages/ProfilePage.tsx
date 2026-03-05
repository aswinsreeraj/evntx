import { useAuthStore } from "../../auth/store/authStore";

function ProfilePage() {
  const { user, roles } = useAuthStore();

  return (
    <div>
      <h2>User Profile</h2>
      <p>Name: {user?.name}</p>
      <p>ID: {user?.id}</p>
      <p>Roles: {roles.join(", ")}</p>
    </div>
  );
}

export default ProfilePage;