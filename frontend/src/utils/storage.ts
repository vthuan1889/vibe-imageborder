// localStorage keys
const STORAGE_KEYS = {
  PRODUCT_FILES: 'vibe_product_files',
  FRAME_FILE: 'vibe_frame_file',
  TEMPLATE_FILE: 'vibe_template_file',
  TEMPLATE_FIELDS: 'vibe_template_fields',
  FIELD_VALUES: 'vibe_field_values',
  OUTPUT_FORMAT: 'vibe_output_format',
  OUTPUT_QUALITY: 'vibe_output_quality',
  OUTPUT_FOLDER: 'vibe_output_folder',
};

export interface AppState {
  productFiles: string[];
  frameFile: string;
  templateFile: string;
  templateFields: string[];
  fieldValues: Record<string, string>;
  format: string;
  quality: number;
  outputFolder: string;
}

export const StorageService = {
  // Save individual values
  saveProductFiles: (files: string[]) => {
    localStorage.setItem(STORAGE_KEYS.PRODUCT_FILES, JSON.stringify(files));
  },

  saveFrameFile: (file: string) => {
    localStorage.setItem(STORAGE_KEYS.FRAME_FILE, file);
  },

  saveTemplateFile: (file: string) => {
    localStorage.setItem(STORAGE_KEYS.TEMPLATE_FILE, file);
  },

  saveTemplateFields: (fields: string[]) => {
    localStorage.setItem(STORAGE_KEYS.TEMPLATE_FIELDS, JSON.stringify(fields));
  },

  saveFieldValues: (values: Record<string, string>) => {
    localStorage.setItem(STORAGE_KEYS.FIELD_VALUES, JSON.stringify(values));
  },

  saveFormat: (format: string) => {
    localStorage.setItem(STORAGE_KEYS.OUTPUT_FORMAT, format);
  },

  saveQuality: (quality: number) => {
    localStorage.setItem(STORAGE_KEYS.OUTPUT_QUALITY, quality.toString());
  },

  saveOutputFolder: (folder: string) => {
    localStorage.setItem(STORAGE_KEYS.OUTPUT_FOLDER, folder);
  },

  // Load individual values
  loadProductFiles: (): string[] => {
    const stored = localStorage.getItem(STORAGE_KEYS.PRODUCT_FILES);
    return stored ? JSON.parse(stored) : [];
  },

  loadFrameFile: (): string => {
    return localStorage.getItem(STORAGE_KEYS.FRAME_FILE) || '';
  },

  loadTemplateFile: (): string => {
    return localStorage.getItem(STORAGE_KEYS.TEMPLATE_FILE) || '';
  },

  loadTemplateFields: (): string[] => {
    const stored = localStorage.getItem(STORAGE_KEYS.TEMPLATE_FIELDS);
    return stored ? JSON.parse(stored) : [];
  },

  loadFieldValues: (): Record<string, string> => {
    const stored = localStorage.getItem(STORAGE_KEYS.FIELD_VALUES);
    return stored ? JSON.parse(stored) : {};
  },

  loadFormat: (): string => {
    return localStorage.getItem(STORAGE_KEYS.OUTPUT_FORMAT) || 'png';
  },

  loadQuality: (): number => {
    const stored = localStorage.getItem(STORAGE_KEYS.OUTPUT_QUALITY);
    return stored ? parseInt(stored) : 90;
  },

  loadOutputFolder: (): string => {
    return localStorage.getItem(STORAGE_KEYS.OUTPUT_FOLDER) || '';
  },

  // Load all state at once
  loadState: (): Partial<AppState> => {
    return {
      productFiles: StorageService.loadProductFiles(),
      frameFile: StorageService.loadFrameFile(),
      templateFile: StorageService.loadTemplateFile(),
      templateFields: StorageService.loadTemplateFields(),
      fieldValues: StorageService.loadFieldValues(),
      format: StorageService.loadFormat(),
      quality: StorageService.loadQuality(),
      outputFolder: StorageService.loadOutputFolder(),
    };
  },

  // Clear all stored data
  clearAll: () => {
    Object.values(STORAGE_KEYS).forEach((key) => {
      localStorage.removeItem(key);
    });
  },
};
