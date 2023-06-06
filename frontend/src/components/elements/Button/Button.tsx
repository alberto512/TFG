'use client';
import styles from './Button.module.css';

type Props = {
  label: string;
  onClick: () => void;
};

const Button = ({ label, onClick }: Props) => {
  return (
    <div className={styles.button} onClick={onClick}>
      <span className={styles.label}>{label}</span>
    </div>
  );
};

export default Button;
