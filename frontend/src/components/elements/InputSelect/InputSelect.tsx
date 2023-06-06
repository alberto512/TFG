import { useEffect, useState } from 'react';
import Image from 'next/image';
import styles from './InputSelect.module.css';

type OptionProps = {
  id: string;
  name: string;
};

type Props = {
  placeholder: string;
  options: OptionProps[];
  onChange: (value: string) => void;
  valueDefault?: string;
  multiple?: boolean;
  onChangeMultiple?: (value: string[]) => void;
};

const InputSelect = ({
  placeholder,
  options,
  onChange,
  valueDefault = '',
  multiple = false,
  onChangeMultiple = () => {},
}: Props) => {
  const [value, setValue] = useState('');
  const [values, setValues] = useState<string[]>([]);
  const [showOptions, setShowOptions] = useState(false);

  const handleBlur = () => {
    setShowOptions(false);
  };

  const handleChange = (option: OptionProps) => {
    if (multiple) {
      if (!values.includes(option.id)) {
        if (values.length === 0) {
          setValue(option.name);
        } else {
          setValue(values.length + 1 + ' options selected');
        }
        setValues([...values, option.id]);
        onChangeMultiple([...values, option.id]);
      } else {
        if (values.length === 2) {
          const filtered = values.filter((value) => value !== option.id)[0];
          setValue(options.find((option) => option.id === filtered)!.name);
        } else if (values.length === 1) {
          setValue(placeholder);
        } else {
          setValue(values.length - 1 + ' options selected');
        }
        setValues([...values.filter((value) => value !== option.id)]);
        onChangeMultiple([...values.filter((value) => value !== option.id)]);
      }
    } else {
      onChange(option.id);
      setValue(option.name);
    }
  };

  useEffect(() => {
    if (valueDefault !== '') {
      setValue(valueDefault);
    } else {
      setValue(placeholder);
    }
  }, [valueDefault]);

  return (
    <div tabIndex={0} className={styles.inputSelect} onClick={() => setShowOptions(!showOptions)} onBlur={handleBlur}>
      <span className={styles.value}>{value}</span>
      <Image className={styles.icon} src={'/svg/select.svg'} alt='Select icon' height={15} width={15} />
      {showOptions && (
        <div className={styles.optionsWrapper}>
          {options.map((option) => (
            <div key={option.id} className={styles.option} onClick={() => handleChange(option)}>
              <span>{option.name}</span>
              {multiple && values.includes(option.id) && (
                <Image src={'/svg/check.svg'} alt='Check icon' height={20} width={20} />
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default InputSelect;
