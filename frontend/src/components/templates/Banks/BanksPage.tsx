'use client';
import { useRouter } from 'next/navigation';
import styles from './BanksPage.module.css';

const BanksPage = () => {
  const router = useRouter();

  const handleSantanderClick = () => {
    router.push('/banks/santander');
  };

  return (
    <div className={styles.banksPage}>
      <span onClick={handleSantanderClick}>Santander</span>
    </div>
  );
};

export default BanksPage;
