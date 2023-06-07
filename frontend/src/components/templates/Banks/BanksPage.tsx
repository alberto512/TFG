'use client';
import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Image from 'next/image';
import styles from './BanksPage.module.css';

const BanksPage = () => {
  const router = useRouter();
  const [isSantanderCompleted, setIsSantanderCompleted] = useState(false);

  const handleSantanderClick = () => {
    router.push('/banks/santander');
  };

  return (
    <div className={styles.banksPage}>
      <div className={styles.cardsContainer}>
        <div className={`${styles.card} ${isSantanderCompleted && styles.cardCompleted}`} onClick={handleSantanderClick}>
          <div className={styles.imageWrapper}>
            <Image src={'/svg/santander.svg'} alt='Santander logo' fill />
          </div>
        </div>
      </div>
    </div>
  );
};

export default BanksPage;
