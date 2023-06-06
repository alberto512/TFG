import Image from 'next/image';
import styles from './HomePage.module.css';

const HomePage = () => {
  return (
    <div className={styles.homePage}>
      <Image src={'/svg/home.svg'} alt='Login image' height={700} width={700} />
    </div>
  );
};

export default HomePage;
