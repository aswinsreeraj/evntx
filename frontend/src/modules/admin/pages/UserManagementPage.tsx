import { useState } from "react";
import { useUsers, useToggleUserStatus } from "../hooks";

function UserManagementPage() {
  const [page, setPage] = useState(1);

  const { data, isLoading } = useUsers({ page, limit: 10 });

  const toggleUser = useToggleUserStatus();

  if (isLoading) return <p>Loading users...</p>;

  return (
    <div>
      <h2>User Management</h2>

      <table>
        <thead>
          <tr>
            <th>Name</th>
            <th>Email</th>
            <th>Status</th>
            <th>Action</th>
          </tr>
        </thead>

        <tbody>
          {data.users.map((user: any) => (
            <tr key={user.user_id}>
              <td>{user.name}</td>
              <td>{user.email}</td>
              <td>{user.is_active ? "Active" : "Blocked"}</td>

              <td>
                <button
                  onClick={() =>
                    toggleUser.mutate({
                      userId: user.user_id,
                      isActive: !user.is_active,
                    })
                  }
                >
                  {user.is_active ? "Block" : "Unblock"}
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

export default UserManagementPage;