import styles from './DotsLoading.module.css';

const DotsLoading = () => {
  return (
    <div className={styles.dotsLoading}>
      <div className={styles.dots}>
        <div></div>
        <div></div>
        <div></div>
      </div>
    </div>
  );
};

export default DotsLoading;
