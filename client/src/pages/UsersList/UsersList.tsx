import { useEffect, useMemo } from "react";
import { useTranslation } from "react-i18next";
import { observer } from "mobx-react-lite";
import { Button, Grid } from "@mui/material";

import { CommonTable, Spinner } from "shared/ui";
import { useDialog } from "shared/hooks";

import { CreateUserForm } from "./ui";
import { UsersModel } from "./model";
import { getColumns } from "./model/columns";

function UsersList() {
  const { t } = useTranslation();

  useEffect(() => {
    UsersModel.fetch();
  }, []);

  const [showCreateUserModal] = useDialog(
    "user:form.createNewUser",
    (hideModal) => <CreateUserForm hideModal={hideModal} />
  );

  const columns = useMemo(() => getColumns(), []);

  return (
    <Grid spacing={2} container direction="column">
      <Grid item alignSelf="flex-end">
        <Button
          variant="contained"
          color="secondary"
          onClick={showCreateUserModal}
        >
          {t("user:form.createNewUser")}
        </Button>
      </Grid>
      <Grid item>
        {UsersModel.loading.has ? (
          <Spinner />
        ) : (
          <CommonTable data={UsersModel.users} columns={columns} />
        )}
      </Grid>
    </Grid>
  );
}

export default observer(UsersList);
