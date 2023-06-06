import Header from '@modules/Header/Header';
import styles from './Default.module.css';

type Props = {
  children: React.ReactNode;
};

const Default = ({ children }: Props) => {
  return (
    <div className={styles.default}>
      <Header />
      <div className={styles.body}>{children}</div>
    </div>
  );
};

export default Default;
