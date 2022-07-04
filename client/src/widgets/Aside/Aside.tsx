import clsx from "clsx";
import { useTranslation } from "react-i18next";
import { observer } from "mobx-react-lite";
import { Drawer, List } from "@mui/material";
import GroupIcon from "@mui/icons-material/Group";

import AsideItemLink from "./AsideItemLink";
import styles from "./Aside.module.sass";

interface AsideProps {
  className?: string;
  isOpen: boolean;
}

const Aside = ({ className, isOpen }: AsideProps) => {
  const { t } = useTranslation();

  return (
    <Drawer
      variant="permanent"
      className={clsx(className, styles.drawer, {
        [styles.drawerOpen]: isOpen,
        [styles.drawerClose]: !isOpen,
      })}
      classes={{
        paper: clsx(styles.paper, styles.drawer, {
          [styles.drawerOpen]: isOpen,
          [styles.drawerClose]: !isOpen,
        }),
      }}
      open={isOpen}
    >
      <List>
        <AsideItemLink
          icon={<GroupIcon />}
          title={t("pages:aside.usersList")}
          to="/usersList"
          isOpen={isOpen}
        />
      </List>
    </Drawer>
  );
};

export default observer(Aside);
