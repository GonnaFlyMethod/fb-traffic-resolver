import { Routes, Route, Navigate } from "react-router-dom";
import { observer } from "mobx-react-lite";

import { UsersList } from "pages";
import { ROUTES } from "shared/consts";

function Router() {
  return (
    <Routes>
      <Route path={ROUTES.USER_LIST} element={<UsersList />} />
      <Route path="/" element={<Navigate to={ROUTES.USER_LIST} replace />} />
      <Route path="*" element={<Navigate to={ROUTES.NOT_FOUND} replace />} />
    </Routes>
  );
}

export default observer(Router);
