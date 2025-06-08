<script lang="ts">
  import DataTable from "./components/DataTable.svelte";
  import Alert from "./components/Alert.svelte";
  import Confirm from "./components/Confirm.svelte";
  import Input from "./components/Input.svelte";
  import WelcomePage from "./components/WelcomePage.svelte";
  import type { TableStats } from "./types/datatable";
  import {
    LoadTable,
    SaveTable,
    AddColumn,
    AddRow,
    AddCalculatedColumn,
    CreateEmptyTable,
    // 基於 ID 的新方法
    LoadTableByID,
    SaveTableByID,
    AddColumnByID,
    AddRowByID,
    AddCalculatedColumnByID,
    CreateEmptyTableByID,
    GetTableCount,
    GetTableInfo,
    GetTableDataByID,
    RemoveTableByID,
    // i18n 方法
    GetText,
    SetLanguage,
    GetCurrentLanguage,
    // 專案檔案操作
    SaveProject,
    SaveProjectAs,
    LoadProject,
    HasUnsavedChanges,
    MarkAsSaved,
    GetCurrentProjectPath,
    // 匯出功能
    ExportTableAsCSV,
    ExportTableAsJSON,
    ExportTableAsExcel,
    // 檔案開啟功能
    OpenCSVFile,
    OpenJSONFile,
    OpenSQLiteFile,
    GetSQLiteTables,
    OpenFileDialog,
  } from "../wailsjs/go/main/App";
  import { onMount } from "svelte";
  import { GetParamValue } from "../wailsjs/go/main/App";
  import {
    showAlert,
    showConfirm,
    showInput,
    alertStore,
    confirmStore,
    inputStore,
    closeAlert,
    closeConfirm,
    closeInput,
  } from "./services/dialogService";

  // 標籤頁介面 - 改為使用數字 ID (slice 索引)
  interface TabInfo {
    id: number;
    name: string;
    isActive: boolean;
  }

  // 狀態管理
  let tabs: TabInfo[] = [{ id: 0, name: "Table 1", isActive: true }];
  let currentTabIndex = 0;
  let isTableLoaded: boolean = false;
  let filePath: string = "";
  let tableKey = 0; // 用於強制重新載入表格組件

  // 歡迎頁面狀態
  let showWelcomePage = true;

  // 標籤頁計數器
  let tabCounter = 1; // 從1開始，因為已有一個 "Table 1"

  // 標籤名稱編輯狀態
  let editingTabIndex: number | null = null;
  let editingTabName = "";
  let editInputRef: HTMLInputElement | null = null;

  // 計算欄輸入狀態（常駐顯示）
  let showColumnInput = true;
  let columnFormulaValue = "";
  let columnNameValue = "";
  let errorMessage = "";
  let showError = false;

  // i18n 狀態
  let currentLanguage = "zh-TW";
  let texts: Record<string, string> = {};

  // 專案狀態
  let currentProjectPath = "";
  let hasUnsavedChanges = false;

  // 統計數據
  let currentStats: TableStats = {
    total_rows: "0",
    total_variables: "0",
    total_cells: "0",
    numeric_variables: "0",
  };

  // i18n 輔助函數
  async function t(key: string): Promise<string> {
    try {
      return await GetText(key);
    } catch (err) {
      console.warn(`翻譯鍵值 "${key}" 不存在，返回預設值`);
      return key;
    }
  }

  // 初始化 i18n
  async function initI18n() {
    try {
      currentLanguage = await GetCurrentLanguage();
      // 預載入常用翻譯
      texts = {
        "ui.buttons.add_variable": await t("ui.buttons.add_variable"),
        "ui.buttons.add_row": await t("ui.buttons.add_row"),
        "ui.buttons.confirm": await t("ui.buttons.confirm"),
        "ui.buttons.clear": await t("ui.buttons.clear"),
        "ui.buttons.save_file": await t("ui.buttons.save_file"),
        "ui.buttons.save_as": await t("ui.buttons.save_as"),
        "ui.buttons.export_table": await t("ui.buttons.export_table"),
        "ui.labels.variable": await t("ui.labels.variable"),
        "ui.labels.row": await t("ui.labels.row"),
        "ui.labels.name": await t("ui.labels.name"),
        "ui.labels.expression": await t("ui.labels.expression"),
        "ui.placeholders.variable_name": await t(
          "ui.placeholders.variable_name"
        ),
        "ui.placeholders.ccl_expression": await t(
          "ui.placeholders.ccl_expression"
        ),
        "ui.placeholders.tab_name": await t("ui.placeholders.tab_name"),
        "ui.stats.total_rows": await t("ui.stats.total_rows"),
        "ui.stats.total_variables": await t("ui.stats.total_variables"),
        "ui.stats.total_cells": await t("ui.stats.total_cells"),
        "ui.stats.numeric_variables": await t("ui.stats.numeric_variables"),
        "ui.stats.basic_statistics": await t("ui.stats.basic_statistics"),
        "ui.defaults.new_variable_name": await t(
          "ui.defaults.new_variable_name"
        ),
        "dialogs.add_variable.title": await t("dialogs.add_variable.title"),
        "dialogs.add_variable.message": await t("dialogs.add_variable.message"),
        "dialogs.add_variable.placeholder": await t(
          "dialogs.add_variable.placeholder"
        ),
        "dialogs.add_variable.confirm": await t("dialogs.add_variable.confirm"),
        "dialogs.add_variable.cancel": await t("dialogs.add_variable.cancel"),
        "dialogs.create_table_failed.title": await t(
          "dialogs.create_table_failed.title"
        ),
        "dialogs.create_table_failed.message": await t(
          "dialogs.create_table_failed.message"
        ),
        "dialog_defaults.alert_title": await t("dialog_defaults.alert_title"),
        "dialog_defaults.confirm_title": await t(
          "dialog_defaults.confirm_title"
        ),
        "dialog_defaults.input_title": await t("dialog_defaults.input_title"),
        "dialog_defaults.confirm_button": await t(
          "dialog_defaults.confirm_button"
        ),
        "dialog_defaults.cancel_button": await t(
          "dialog_defaults.cancel_button"
        ),
      };
    } catch (err) {
      console.error("i18n 初始化失敗:", err);
    }
  }

  // 組件掛載時執行
  onMount(async () => {
    // 初始化 i18n
    await initI18n();

    // 為初始標籤頁創建空白資料表
    try {
      const initialTab = tabs[0];
      if (initialTab) {
        const actualTableID = await CreateEmptyTableByID(
          initialTab.id,
          initialTab.name
        );
        if (actualTableID >= 0) {
          tabs[0].id = actualTableID;
          isTableLoaded = true;
          tableKey++; // 觸發表格重新載入
          console.log(
            `為初始標籤頁 ${initialTab.name} 創建空白資料表，ID: ${actualTableID}`
          );
        } else {
          console.warn("無法為初始標籤頁創建資料表");
        }
      }
    } catch (err) {
      console.error("初始化標籤頁資料表時發生錯誤:", err);
    }

    // 獲取命令行傳入的檔案路徑
    try {
      const autoLoadPath = (await GetParamValue("filepath")) || "";

      if (autoLoadPath) {
        filePath = autoLoadPath;
        // 如果提供了文件路徑，則自動載入
        await handleLoadTable();
      }
    } catch (err) {
      console.error("無法獲取啟動參數", err);
    }
  });

  // 標籤頁操作
  async function addNewTab() {
    tabCounter++; // 增加計數器
    const newTabName = `Table ${tabCounter}`;
    const newTabID = tabs.length; // 使用數字ID作為slice索引

    // 為新標籤頁創建空白資料表
    try {
      const actualTableID = await CreateEmptyTableByID(newTabID, newTabName);
      if (actualTableID >= 0) {
        // CreateEmptyTableByID 返回 number (tableID)，-1表示失敗
        const newTab: TabInfo = {
          id: actualTableID, // 使用實際返回的 table ID
          name: newTabName,
          isActive: false,
        };

        // 設置所有標籤為非活動
        tabs = tabs.map((tab) => ({ ...tab, isActive: false }));
        // 添加新標籤並設為活動
        tabs = [...tabs, { ...newTab, isActive: true }];
        currentTabIndex = tabs.length - 1;

        isTableLoaded = true;

        // 強制重新載入表格組件
        tableKey++;

        console.log(
          `成功為標籤頁 ${newTabName} 創建空白資料表，ID: ${actualTableID}`
        );
      } else {
        console.error(`為標籤頁 ${newTabName} 創建空白資料表失敗`);
        await showAlert({
          title: "創建失敗",
          message: "創建新標籤頁失敗",
          type: "error",
        });
      }
    } catch (err) {
      console.error("創建空白資料表時發生錯誤:", err);
      await showAlert({
        title: "創建錯誤",
        message: `創建新標籤頁時發生錯誤: ${err}`,
        type: "error",
      });
    }
  }

  async function switchTab(index: number) {
    // 更新標籤頁狀態
    tabs = tabs.map((tab, i) => ({ ...tab, isActive: i === index }));
    currentTabIndex = index;

    // 檢查切換到的標籤頁是否有有效的資料表
    const currentTab = tabs[index];
    if (currentTab && currentTab.id >= 0) {
      try {
        // 嘗試獲取表格資料以驗證是否存在
        const data = await GetTableDataByID(currentTab.id);
        if (data && (data.rows || data.columns)) {
          isTableLoaded = true;
        } else {
          isTableLoaded = false;
        }
      } catch (err) {
        console.log(`標籤頁 ${index} 的資料表不存在或無效:`, err);
        isTableLoaded = false;
      }
    } else {
      isTableLoaded = false;
    }

    // 強制重新載入表格組件
    tableKey++;
  }

  // 刪除標籤頁
  async function removeTab(index: number, event?: Event) {
    // 阻止事件冒泡，避免觸發 switchTab
    if (event) {
      event.stopPropagation();
    }

    // 至少保留一個標籤頁
    if (tabs.length <= 1) {
      await showAlert({
        title: "無法刪除",
        message: "至少需要保留一個標籤頁",
        type: "warning",
      });
      return;
    }

    const tabToRemove = tabs[index];
    const confirmResult = await showConfirm({
      title: "確認刪除",
      message: `確定要刪除標籤頁 "${tabToRemove.name}" 嗎？此操作無法復原。`,
      type: "danger",
      confirmText: "刪除",
      cancelText: "取消",
    });

    if (!confirmResult) {
      return;
    }

    try {
      // 調用後端API刪除表格
      const success = await RemoveTableByID(tabToRemove.id);
      if (!success) {
        console.warn(`刪除表格 ID ${tabToRemove.id} 失敗，但仍會移除標籤頁`);
      }

      // 從tabs陣列中移除對應的tab
      tabs = tabs.filter((_, i) => i !== index);

      // 處理刪除後的tab切換邏輯
      if (index === currentTabIndex) {
        // 如果刪除的是當前活動標籤頁
        if (index >= tabs.length) {
          // 如果刪除的是最後一個標籤頁，切換到前一個
          currentTabIndex = tabs.length - 1;
        }
        // 否則保持當前索引（會自動切換到下一個標籤頁）

        // 設置新的活動標籤頁
        if (tabs.length > 0) {
          tabs = tabs.map((tab, i) => ({
            ...tab,
            isActive: i === currentTabIndex,
          }));

          // 檢查新活動標籤頁的資料表狀態
          const newActiveTab = tabs[currentTabIndex];
          if (newActiveTab && newActiveTab.id >= 0) {
            try {
              const data = await GetTableDataByID(newActiveTab.id);
              isTableLoaded = !!(data && (data.rows || data.columns));
            } catch (err) {
              console.log(
                `標籤頁 ${currentTabIndex} 的資料表不存在或無效:`,
                err
              );
              isTableLoaded = false;
            }
          } else {
            isTableLoaded = false;
          }

          // 強制重新載入表格組件
          tableKey++;
        }
      } else if (index < currentTabIndex) {
        // 如果刪除的標籤頁在當前活動標籤頁之前，需要調整當前索引
        currentTabIndex--;
      }

      console.log(`成功刪除標籤頁 "${tabToRemove.name}"`);
    } catch (err) {
      console.error("刪除標籤頁時發生錯誤:", err);
      await showAlert({
        title: "刪除錯誤",
        message: `刪除標籤頁時發生錯誤: ${err}`,
        type: "error",
      });
    }
  }

  // 功能列操作
  async function addColumn() {
    // 檢查是否有活動的資料表
    if (!isTableLoaded) {
      // 如果沒有資料表，先創建一個空白資料表
      const activeTableID = tabs[currentTabIndex]?.id ?? 0;
      const createSuccess = await CreateEmptyTableByID(
        activeTableID,
        `Table ${activeTableID + 1}`
      );
      if (createSuccess >= 0) {
        isTableLoaded = true;
        // 更新標籤頁 ID 為實際的 table ID
        tabs[currentTabIndex].id = createSuccess;
      } else {
        await showAlert({
          title: texts["dialogs.create_table_failed.title"] || "創建失敗",
          message:
            texts["dialogs.create_table_failed.message"] || "無法創建資料表",
          type: "error",
        });
        return;
      }
    }

    const columnName = await showInput({
      title: texts["dialogs.add_variable.title"] || "新增變項",
      message: texts["dialogs.add_variable.message"] || "請輸入新變項名稱:",
      placeholder: texts["dialogs.add_variable.placeholder"] || "變項名稱",
      defaultValue: `${texts["ui.defaults.new_variable_name"] || "新變項"} ${currentStats["total_variables"] ? parseInt(currentStats["total_variables"]) + 1 : 1}`,
      confirmText: texts["dialogs.add_variable.confirm"] || "新增",
      cancelText: texts["dialogs.add_variable.cancel"] || "取消",
    });
    console.log("showInput 返回值:", columnName);
    if (columnName) {
      const activeTableID = tabs[currentTabIndex]?.id ?? 0;
      console.log("正在調用 AddColumnByID，參數:", {
        activeTableID,
        columnName,
      });
      try {
        const success = await AddColumnByID(activeTableID, columnName);
        console.log("AddColumnByID 回傳結果:", success);
        if (success) {
          // 重新載入表格數據以顯示新增的欄位
          await refreshCurrentTable();
          console.log("新增欄位成功");
        } else {
          console.error("AddColumn 回傳 false");
          await showAlert({
            title: "新增失敗",
            message: "新增欄位失敗",
            type: "error",
          });
        }
      } catch (error) {
        console.error("AddColumn 發生錯誤:", error);
        await showAlert({
          title: "新增錯誤",
          message: `新增欄位發生錯誤: ${error}`,
          type: "error",
        });
      }
    }
  }

  async function addRow() {
    // 檢查是否有活動的資料表
    if (!isTableLoaded) {
      // 如果沒有資料表，先創建一個空白資料表
      const activeTableID = tabs[currentTabIndex]?.id ?? 0;
      const createSuccess = await CreateEmptyTableByID(
        activeTableID,
        `Table ${activeTableID + 1}`
      );
      if (createSuccess >= 0) {
        isTableLoaded = true;
        // 更新標籤頁 ID 為實際的 table ID
        tabs[currentTabIndex].id = createSuccess;
      } else {
        await showAlert({
          title: "創建失敗",
          message: "無法創建資料表",
          type: "error",
        });
        return;
      }
    }

    const activeTableID = tabs[currentTabIndex]?.id ?? 0;
    console.log("正在調用 AddRowByID，參數:", { activeTableID });
    try {
      const success = await AddRowByID(activeTableID);
      console.log("AddRowByID 回傳結果:", success);
      if (success) {
        // 重新載入表格數據以顯示新增的行
        await refreshCurrentTable();
        console.log("新增行成功");
      } else {
        console.error("AddRowByID 回傳 false");
        await showAlert({
          title: "新增失敗",
          message: "新增行失敗",
          type: "error",
        });
      }
    } catch (error) {
      console.error("AddRow 發生錯誤:", error);
      await showAlert({
        title: "新增錯誤",
        message: `新增行發生錯誤: ${error}`,
        type: "error",
      });
    }
  }

  // 移除了 addCalculatedColumn 函數，因為輸入框常駐顯示

  async function confirmAddColumn() {
    if (columnNameValue && columnFormulaValue) {
      const activeTableID = tabs[currentTabIndex]?.id ?? 0;
      const success = await AddCalculatedColumnByID(
        activeTableID,
        columnNameValue,
        columnFormulaValue
      );
      if (success) {
        // 重新載入表格數據以顯示新增的計算欄
        await refreshCurrentTable();
        clearColumnInput();
      } else {
        showError = true;
        errorMessage = "添加計算欄失敗";
      }
    } else {
      showError = true;
      errorMessage = "請輸入欄位名稱與 CCL 表達式";
    }
  }

  function clearColumnInput() {
    columnFormulaValue = "";
    columnNameValue = "";
    showError = false;
    errorMessage = "";
  }

  // 重新載入當前表格
  async function refreshCurrentTable() {
    // 通過改變 key 來強制重新載入表格組件
    tableKey++;
  }

  // 接收統計数據更新
  function handleStatsUpdate(event: CustomEvent) {
    currentStats = event.detail;
  }

  // 底部工具列操作
  async function openFile() {
    // 檢查是否有未儲存的變更
    if (await HasUnsavedChanges()) {
      const confirmed = await showConfirm({
        title: await t("file_operations.unsaved_changes"),
        message: await t("file_operations.save_before_close"),
        confirmText: await t("ui.buttons.confirm"),
        cancelText: await t("ui.buttons.cancel"),
        type: "warning",
      });

      if (confirmed) {
        await saveProject();
      } else {
        // 如果用戶取消儲存，則不繼續開啟檔案
        return;
      }
    }

    // 開啟專案檔案
    try {
      // Wails OpenFileDialog 預期一個篩選器字串，例如 "Insyra 專案檔案 (*.insa)|*.insa"
      const fileFilter = `${await t("file_operations.project_file_description")} (*.insa)|*.insa`;
      const selectedPath = await OpenFileDialog(fileFilter);

      if (selectedPath) {
        const success = await LoadProject(selectedPath);
        if (success) {
          currentProjectPath = selectedPath;
          hasUnsavedChanges = false;
          await showAlert({
            title: await t("messages.import_success"),
            message: selectedPath,
            type: "success",
          });
          location.reload();
        } else {
          await showAlert({
            title: await t("messages.import_fail"),
            message: selectedPath,
            type: "error",
          });
        }
      }
    } catch (err) {
      console.error("開啟專案檔案時發生錯誤:", err);
      await showAlert({
        title: await t("messages.import_fail"),
        message: String(err),
        type: "error",
      });
    }
  }

  async function saveProject() {
    if (!currentProjectPath) {
      await saveProjectAs();
      return;
    }

    const success = await SaveProject(currentProjectPath);
    if (success) {
      hasUnsavedChanges = false;
      await MarkAsSaved();
      await showAlert({
        title: await t("messages.save_success"),
        message: currentProjectPath,
        type: "success",
      });
    } else {
      await showAlert({
        title: await t("messages.save_fail"),
        message: currentProjectPath,
        type: "error",
      });
    }
  }

  async function saveProjectAs() {
    try {
      const fileFilter = `${await t("file_operations.project_file_description")} (*.insa)|*.insa`;
      const selectedPath = await OpenFileDialog(fileFilter); // 實際上 Wails 的 OpenFileDialog 通常用於開啟，另存為可能需要 SaveFileDialog
      // 假設此處的 OpenFileDialog 也能用於「另存為」的情境，或者後端 SaveProjectAs 會處理路徑

      if (selectedPath) {
        const success = await SaveProjectAs(selectedPath);
        if (success) {
          currentProjectPath = selectedPath;
          hasUnsavedChanges = false;
          await MarkAsSaved();
          await showAlert({
            title: await t("messages.save_success"),
            message: selectedPath,
            type: "success",
          });
        } else {
          await showAlert({
            title: await t("messages.save_fail"),
            message: selectedPath,
            type: "error",
          });
        }
      }
    } catch (err) {
      console.error("另存專案檔案時發生錯誤:", err);
      await showAlert({
        title: await t("messages.save_fail"),
        message: String(err),
        type: "error",
      });
    }
  }

  async function exportCurrentTable() {
    if (!isTableLoaded) {
      await showAlert({
        title: await t("messages.export_fail"),
        message: "沒有可匯出的資料表",
        type: "warning",
      });
      return;
    }

    // 選擇匯出格式
    const format = await showInput({
      title: await t("dialogs.export.title"),
      message: await t("dialogs.export.message"),
      placeholder: "csv, json, excel",
      defaultValue: "csv",
      confirmText: await t("ui.buttons.export_table"),
      cancelText: await t("ui.buttons.cancel"),
    });

    if (!format) return;

    // 選擇匯出路徑
    const currentTabName = tabs[currentTabIndex]?.name || "table";
    const defaultFileName = `${currentTabName}.${format.toLowerCase()}`;
    let fileFilter = "";
    switch (format.toLowerCase()) {
      case "csv":
        fileFilter = `CSV (*.csv)|*.csv`;
        break;
      case "json":
        fileFilter = `JSON (*.json)|*.json`;
        break;
      case "excel":
      case "xlsx":
        fileFilter = `Excel (*.xlsx)|*.xlsx`;
        break;
      default:
        await showAlert({
          title: await t("messages.export_fail"),
          message: "不支援的匯出格式",
          type: "error",
        });
        return;
    }

    try {
      const selectedPath = await OpenFileDialog(fileFilter); // 假設 OpenFileDialog 可用於儲存

      if (!selectedPath) return;

      const currentTableID = tabs[currentTabIndex]?.id ?? 0;
      let success = false;

      switch (format.toLowerCase()) {
        case "csv":
          success = await ExportTableAsCSV(currentTableID, selectedPath);
          break;
        case "json":
          success = await ExportTableAsJSON(currentTableID, selectedPath);
          break;
        case "excel":
        case "xlsx":
          success = await ExportTableAsExcel(currentTableID, selectedPath);
          break;
      }

      if (success) {
        await showAlert({
          title: await t("messages.export_success"),
          message: selectedPath,
          type: "success",
        });
      } else {
        await showAlert({
          title: await t("messages.export_fail"),
          message: selectedPath,
          type: "error",
        });
      }
    } catch (err) {
      console.error("匯出表格時發生錯誤:", err);
      await showAlert({
        title: await t("messages.export_fail"),
        message: String(err),
        type: "error",
      });
    }
  }

  async function openSettings() {
    // TODO: 實現設定功能，可以設置表格外觀、預設值等
    console.log("開啟設定");
    await showAlert({
      title: "功能開發中",
      message: "設定功能尚未實現",
      type: "info",
    });
  }

  // 載入資料表
  async function handleLoadTable() {
    // filePath 變數在此函數中不再由用戶手動輸入，而是來自例如啟動參數
    // 因此，如果 filePath 為空，表示沒有透過參數傳入檔案路徑，
    // 此函數可能不應該被直接調用，或者應該有其他方式獲取路徑。
    // 這裡假設如果 filePath 為空，則不執行任何操作，或提示用戶。
    if (!filePath) {
      // 考慮是否應該提示用戶選擇檔案，如果此函數預期被用戶觸發
      // 例如，可以調用一個新的輔助函數來開啟檔案選擇器
      // await selectAndLoadTable();
      // 或者，如果 filePath 必須由外部設定，則顯示警告
      await showAlert({
        title: "載入錯誤",
        message: "沒有指定檔案路徑", // 更新訊息以反映 filePath 的來源
        type: "warning",
      });
      return;
    }

    try {
      const activeTableID = tabs[currentTabIndex]?.id ?? 0;
      const tableName = `Table ${activeTableID + 1}`;
      const newTableID = await LoadTableByID(
        activeTableID,
        tableName,
        filePath
      );
      if (newTableID >= 0) {
        isTableLoaded = true;
        // 更新標籤頁 ID 為實際的 table ID
        tabs[currentTabIndex].id = newTableID;
      } else {
        await showAlert({
          title: "載入失敗",
          message: "載入資料表失敗",
          type: "error",
        });
      }
    } catch (err) {
      await showAlert({
        title: "載入錯誤",
        message: `發生錯誤: ${err}`,
        type: "error",
      });
    }
  }

  // 儲存資料表
  async function handleSaveTable() {
    if (!isTableLoaded || !filePath) {
      await showAlert({
        title: "儲存錯誤",
        message: "請先載入資料表或指定儲存路徑",
        type: "warning",
      });
      return;
    }

    try {
      const activeTableID = tabs[currentTabIndex]?.id ?? 0;
      const success = await SaveTableByID(activeTableID, filePath);
      if (success) {
        await showAlert({
          title: "儲存成功",
          message: "資料表已成功儲存",
          type: "success",
        });
      } else {
        await showAlert({
          title: "儲存失敗",
          message: "儲存資料表失敗",
          type: "error",
        });
      }
    } catch (err) {
      await showAlert({
        title: "儲存錯誤",
        message: `發生錯誤: ${err}`,
        type: "error",
      });
    }
  }

  // 標籤名稱編輯功能
  function startEditingTabName(index: number, event?: Event) {
    if (event) {
      event.stopPropagation();
    }
    editingTabIndex = index;
    editingTabName = tabs[index].name;

    // 使用 setTimeout 確保 DOM 已更新
    setTimeout(() => {
      if (editInputRef) {
        editInputRef.focus();
        editInputRef.select();
      }
    }, 0);
  }

  function finishEditingTabName() {
    if (editingTabIndex !== null && editingTabName.trim()) {
      tabs[editingTabIndex].name = editingTabName.trim();
      tabs = [...tabs]; // 觸發重新渲染
    }
    editingTabIndex = null;
    editingTabName = "";
  }

  function cancelEditingTabName() {
    editingTabIndex = null;
    editingTabName = "";
  }

  function handleTabNameKeydown(event: KeyboardEvent) {
    if (event.key === "Enter") {
      finishEditingTabName();
    } else if (event.key === "Escape") {
      cancelEditingTabName();
    }
  }

  function handleTabDoubleClick(index: number, event: Event) {
    event.preventDefault();
    event.stopPropagation();
    startEditingTabName(index);
  }

  function handleTabRightClick(index: number, event: MouseEvent) {
    event.preventDefault();
    event.stopPropagation();
    startEditingTabName(index);
  }

  // 歡迎頁面事件處理
  async function handleWelcomeAction(event: CustomEvent) {
    const { type } = event.detail;

    try {
      switch (type) {
        case "open_csv":
          await handleOpenCSV();
          break;
        case "open_json":
          await handleOpenJSON();
          break;
        case "open_sqlite":
          await handleOpenSQLite();
          break;
        case "open_project":
          await handleOpenProject();
          break;
        case "new_project":
          await handleNewProject();
          break;
        default:
          console.warn("未知的歡迎頁面操作:", type);
      }
    } catch (err) {
      console.error("處理歡迎頁面操作時發生錯誤:", err);
      await showAlert({
        title: "操作失敗",
        message: `執行操作時發生錯誤: ${err}`,
        type: "error",
      });
    }
  }

  // 開啟 CSV 檔案
  async function handleOpenCSV() {
    try {
      const filePath = await OpenFileDialog("CSV 檔案 (*.csv)|*.csv");
      if (filePath) {
        const tableId = await OpenCSVFile(filePath);
        if (tableId >= 0) {
          // 成功開啟，隱藏歡迎頁面
          showWelcomePage = false;
          // 創建新標籤頁
          await createTabFromFile(filePath, tableId, "csv");
        } else {
          await showAlert({
            title: "開啟失敗",
            message: "無法開啟 CSV 檔案",
            type: "error",
          });
        }
      }
    } catch (err) {
      console.error("開啟 CSV 檔案失敗:", err);
      await showAlert({
        title: "開啟錯誤",
        message: `開啟 CSV 檔案時發生錯誤: ${err}`,
        type: "error",
      });
    }
  }

  // 開啟 JSON 檔案
  async function handleOpenJSON() {
    try {
      const filePath = await OpenFileDialog("JSON 檔案 (*.json)|*.json");
      if (filePath) {
        const tableId = await OpenJSONFile(filePath);
        if (tableId >= 0) {
          // 成功開啟，隱藏歡迎頁面
          showWelcomePage = false;
          // 創建新標籤頁
          await createTabFromFile(filePath, tableId, "json");
        } else {
          await showAlert({
            title: "開啟失敗",
            message: "無法開啟 JSON 檔案",
            type: "error",
          });
        }
      }
    } catch (err) {
      console.error("開啟 JSON 檔案失敗:", err);
      await showAlert({
        title: "開啟錯誤",
        message: `開啟 JSON 檔案時發生錯誤: ${err}`,
        type: "error",
      });
    }
  }

  // 開啟 SQLite 檔案
  async function handleOpenSQLite() {
    try {
      const filePath = await OpenFileDialog(
        "SQLite 檔案 (*.db;*.sqlite;*.sqlite3)|*.db;*.sqlite;*.sqlite3"
      );
      if (filePath) {
        // 首先獲取表格列表
        const tables = await GetSQLiteTables(filePath);
        if (tables && tables.length > 0) {
          // 如果有多個表格，讓用戶選擇
          let selectedTable = tables[0]; // 預設選擇第一個

          if (tables.length > 1) {
            // 顯示表格選擇對話框
            const tableList = tables.join("\n");
            const selected = await showInput({
              title: "選擇資料表",
              message: `發現多個資料表，請輸入要開啟的表格名稱：\n\n可用的表格：\n${tableList}`,
              placeholder: tables[0],
              defaultValue: tables[0],
            });

            if (selected && tables.includes(selected)) {
              selectedTable = selected;
            } else if (!selected) {
              return; // 用戶取消
            }
          }

          const tableId = await OpenSQLiteFile(filePath, selectedTable);
          if (tableId >= 0) {
            // 成功開啟，隱藏歡迎頁面
            showWelcomePage = false;
            // 創建新標籤頁
            await createTabFromFile(filePath, tableId, "sqlite", selectedTable);
          } else {
            await showAlert({
              title: "開啟失敗",
              message: `無法開啟 SQLite 資料表: ${selectedTable}`,
              type: "error",
            });
          }
        } else {
          await showAlert({
            title: "開啟失敗",
            message: "此 SQLite 檔案中沒有找到資料表",
            type: "error",
          });
        }
      }
    } catch (err) {
      console.error("開啟 SQLite 檔案失敗:", err);
      await showAlert({
        title: "開啟錯誤",
        message: `開啟 SQLite 檔案時發生錯誤: ${err}`,
        type: "error",
      });
    }
  }

  // 開啟專案檔案
  async function handleOpenProject() {
    try {
      const filePath = await OpenFileDialog("Insyra 專案檔案 (*.insa)|*.insa");
      if (filePath) {
        const success = await LoadProject(filePath);
        if (success) {
          // 成功開啟，隱藏歡迎頁面
          showWelcomePage = false;
          // 重新載入所有標籤頁
          await refreshAllTabs();
        } else {
          await showAlert({
            title: "開啟失敗",
            message: "無法開啟專案檔案",
            type: "error",
          });
        }
      }
    } catch (err) {
      console.error("開啟專案檔案失敗:", err);
      await showAlert({
        title: "開啟錯誤",
        message: `開啟 專案檔案時發生錯誤: ${err}`,
        type: "error",
      });
    }
  }

  // 建立新專案
  async function handleNewProject() {
    console.log("處理建立新專案");
    // 隱藏歡迎頁面，進入空白專案
    showWelcomePage = false;
    console.log("歡迎頁面已關閉");

    // 重設所有狀態
    tabs = [{ id: 0, name: "Table 1", isActive: true }];
    currentTabIndex = 0;
    tabCounter = 1;

    // 創建空白資料表
    try {
      const actualTableID = await CreateEmptyTableByID(0, "Table 1");
      if (actualTableID >= 0) {
        tabs[0].id = actualTableID;
        isTableLoaded = true;
        tableKey++;
        console.log(`為新專案創建空白資料表，ID: ${actualTableID}`);
      } else {
        console.warn("為新專案創建空白資料表失敗");
      }
    } catch (err) {
      console.error("創建新專案時發生錯誤:", err);
    }
  }

  // 從檔案創建標籤頁的輔助函數
  async function createTabFromFile(
    filePath: string,
    tableId: number,
    fileType: string,
    tableName?: string
  ) {
    const fileName = filePath.split("\\").pop()?.split("/").pop() || "Unknown";
    const tabName = tableName ? `${fileName} - ${tableName}` : fileName;

    // 清空現有標籤頁並創建新的標籤頁
    tabs = [{ id: tableId, name: tabName, isActive: true }];
    currentTabIndex = 0;
    tabCounter = 1;
    isTableLoaded = true;
    tableKey++;

    console.log(
      `從 ${fileType.toUpperCase()} 檔案創建標籤頁: ${tabName}, ID: ${tableId}`
    );
  }

  // 重新載入所有標籤頁的輔助函數
  async function refreshAllTabs() {
    try {
      const tableCount = await GetTableCount();
      tabs = [];

      for (let i = 0; i < tableCount; i++) {
        const tableInfo = await GetTableInfo(i);
        if (tableInfo) {
          tabs.push({
            id: i,
            name: tableInfo.name || `Table ${i + 1}`,
            isActive: i === 0,
          });
        }
      }

      if (tabs.length > 0) {
        currentTabIndex = 0;
        isTableLoaded = true;
        tableKey++;
      } else {
        // 如果沒有標籤頁，創建一個空白標籤頁
        await handleNewProject();
      }
    } catch (err) {
      console.error("重新載入標籤頁失敗:", err);
      // 創建一個空白標籤頁作為備用
      await handleNewProject();
    }
  }

  // ...existing code...
</script>

<main>
  <!-- 歡迎頁面 -->
  {#if showWelcomePage}
    <WelcomePage on:action={handleWelcomeAction} />
  {:else}
    <!-- 標籤列 -->
    <div class="tab-bar">
      <div class="tab-row">
        {#each tabs as tab, index}
          <div class="tab-container">
            <button
              class="tab-button"
              class:tab-active={tab.isActive}
              on:click={() => switchTab(index)}
              on:dblclick={(event) => handleTabDoubleClick(index, event)}
              on:contextmenu={(event) => handleTabRightClick(index, event)}
            >
              {#if editingTabIndex === index}
                <input
                  bind:this={editInputRef}
                  type="text"
                  class="tab-name-input"
                  bind:value={editingTabName}
                  on:blur={finishEditingTabName}
                  on:keydown={handleTabNameKeydown}
                  placeholder={texts["ui.placeholders.tab_name"] || "標籤名稱"}
                />
              {:else}
                {tab.name}
              {/if}
            </button>
            <button
              class="tab-close-button"
              class:disabled={tabs.length <= 1}
              on:click={(event) => removeTab(index, event)}
              title="刪除標籤頁"
            >
              ×
            </button>
          </div>
        {/each}
        <button class="tab-add-button" on:click={addNewTab} tabindex="-1"
          >+</button
        >
      </div>
    </div>

    <!-- 功能列 -->
    <div class="function-bar">
      <div class="function-buttons">
        <button class="function-button" on:click={addColumn}>
          {texts["ui.buttons.add_variable"] || "新增變項"}
        </button>
        <button class="function-button" on:click={addRow}>
          {texts["ui.buttons.add_row"] || "新增列"}
        </button>
      </div>

      <!-- 計算變項輸入區域（常駐顯示） -->
      <div class="column-input">
        <div class="input-row">
          <span class="fx-label">fx</span>
          <input
            type="text"
            class="column-name-input"
            placeholder={texts["ui.placeholders.variable_name"] || "變項名稱"}
            bind:value={columnNameValue}
          />
          <span class="equals">=</span>
          <input
            type="text"
            class="formula-input"
            placeholder={texts["ui.placeholders.ccl_expression"] ||
              "CCL 表達式"}
            bind:value={columnFormulaValue}
          />
          <button class="confirm-button" on:click={confirmAddColumn}>
            {texts["ui.buttons.confirm"] || "✓"}
          </button>
          <button class="cancel-button" on:click={clearColumnInput}>
            {texts["ui.buttons.clear"] || "清除"}
          </button>
        </div>
        {#if showError}
          <div class="error-message">{errorMessage}</div>
        {/if}
      </div>
    </div>

    <!-- 主要內容區域 -->
    <div class="main-content">
      <!-- 左側表格區域 -->
      <div class="table-area">
        {#if isTableLoaded}
          <DataTable
            tableID={tabs[currentTabIndex]?.id ?? 0}
            on:statsUpdate={handleStatsUpdate}
            {tableKey}
          />
        {:else}
          <div class="table-placeholder">
            <p>
              {texts["ui.table.table_placeholder"] || "資料表為空，請新增資料"}
            </p>
          </div>
        {/if}
      </div>

      <!-- 右側資訊區域 -->
      <div class="info-area">
        <div class="info-header">
          <h3>{texts["ui.stats.basic_statistics"] || "基本統計"}</h3>
        </div>
        <div class="stats-content">
          <div class="stat-item">
            <span class="stat-label"
              >{texts["ui.stats.total_rows"] || "總列數"}:</span
            >
            <span class="stat-value">{currentStats["total_rows"]}</span>
          </div>
          <div class="stat-item">
            <span class="stat-label"
              >{texts["ui.stats.total_variables"] || "總變項數"}:</span
            >
            <span class="stat-value">{currentStats["total_variables"]}</span>
          </div>
          <div class="stat-item">
            <span class="stat-label"
              >{texts["ui.stats.total_cells"] || "總儲存格"}:</span
            >
            <span class="stat-value">{currentStats["total_cells"]}</span>
          </div>
          <div class="stat-item">
            <span class="stat-label"
              >{texts["ui.stats.numeric_variables"] || "數值變項數"}:</span
            >
            <span class="stat-value">{currentStats["numeric_variables"]}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 底部工具列 -->
    <div class="bottom-toolbar">
      <button class="toolbar-button open-button" on:click={openFile}>
        {texts["ui.buttons.open_file"] || "開啟專案"}
      </button>
      <button class="toolbar-button save-button" on:click={saveProject}>
        {texts["ui.buttons.save_file"] || "儲存專案"}
      </button>
      <button class="toolbar-button save-as-button" on:click={saveProjectAs}>
        {texts["ui.buttons.save_as"] || "另存新檔"}
      </button>
      <button
        class="toolbar-button export-button"
        on:click={exportCurrentTable}
      >
        {texts["ui.buttons.export_table"] || "匯出資料表"}
      </button>
    </div>
  {/if}
</main>

<!-- 對話框組件 -->
<Alert
  visible={$alertStore.visible}
  options={$alertStore.options}
  {texts}
  on:close={(e) => closeAlert()}
/>

<Confirm
  visible={$confirmStore.visible}
  options={$confirmStore.options}
  {texts}
  on:close={(e) => closeConfirm(e.detail.result)}
/>

<Input
  visible={$inputStore.visible}
  options={$inputStore.options}
  {texts}
  on:close={(e) => {
    console.log("Input 對話框關閉事件:", e.detail);
    closeInput(e.detail.result);
  }}
/>

<style>
  :global(body) {
    margin: 0;
    padding: 0;
    font-family:
      "Nunito",
      -apple-system,
      BlinkMacSystemFont,
      "Segoe UI",
      Roboto,
      Oxygen,
      Ubuntu,
      Cantarell,
      "Open Sans",
      "Helvetica Neue",
      sans-serif;
    background: var(--background-color);
    overflow: hidden;
  }

  main {
    display: flex;
    flex-direction: column;
    height: 100vh;
    width: 100%;
    background: var(--background-color);
    position: relative;
    overflow: hidden;
  }

  /* 背景裝飾 */
  main::before {
    content: "";
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-image: radial-gradient(
        circle at 20% 20%,
        rgba(255, 255, 255, 0.1) 1px,
        transparent 1px
      ),
      radial-gradient(
        circle at 80% 80%,
        rgba(255, 255, 255, 0.1) 1px,
        transparent 1px
      );
    background-size:
      100px 100px,
      120px 120px;
    pointer-events: none;
    z-index: 0;
  }

  /* 標籤列樣式 */
  .tab-bar {
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(20px);
    border-bottom: 1px solid rgba(255, 255, 255, 0.2);
    padding: var(--spacing-sm) var(--spacing-md) 0;
    box-shadow: var(--shadow-1);
    position: relative;
    z-index: 10;
  }

  .tab-row {
    display: flex;
    flex-wrap: wrap;
    gap: var(--spacing-xs);
    align-items: center;
  }

  .tab-container {
    position: relative;
    display: flex;
    align-items: center;
    border-radius: var(--radius-medium) var(--radius-medium) 0 0;
    overflow: hidden;
  }

  .tab-button {
    padding: var(--spacing-sm) var(--spacing-lg);
    border: none;
    border-radius: var(--radius-medium) var(--radius-medium) 0 0;
    font-size: 0.9rem;
    font-weight: 500;
    cursor: pointer;
    transition: all var(--transition-standard);
    background: linear-gradient(
      135deg,
      rgba(255, 255, 255, 0.7),
      rgba(248, 250, 252, 0.8)
    );
    color: var(--text-secondary);
    margin-bottom: -1px;
    padding-right: var(--spacing-xl);
    position: relative;
    backdrop-filter: blur(10px);
    border: 1px solid rgba(255, 255, 255, 0.3);
    border-bottom: none;
    box-shadow: var(--shadow-1);
  }

  .tab-button.tab-active {
    background: linear-gradient(
      135deg,
      rgba(255, 255, 255, 0.95),
      rgba(248, 250, 252, 0.95)
    );
    color: var(--text-primary);
    box-shadow: var(--shadow-2);
    transform: translateY(-2px);
    border-color: rgba(255, 255, 255, 0.5);
  }

  .tab-button::before {
    content: "";
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 2px;
    background: linear-gradient(
      90deg,
      transparent,
      var(--primary-color),
      transparent
    );
    opacity: 0;
    transition: opacity var(--transition-standard);
  }

  .tab-button.tab-active::before {
    opacity: 1;
  }

  .tab-button:hover {
    background: linear-gradient(
      135deg,
      rgba(255, 255, 255, 0.9),
      rgba(248, 250, 252, 0.9)
    );
    transform: translateY(-1px);
    box-shadow: var(--shadow-2);
  }

  .tab-close-button {
    position: absolute;
    right: var(--spacing-xs);
    top: 50%;
    transform: translateY(-50%);
    width: 24px;
    height: 24px;
    border: none;
    border-radius: 50%;
    background: rgba(0, 0, 0, 0.05);
    color: var(--text-secondary);
    font-size: 14px;
    font-weight: bold;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all var(--transition-fast);
    opacity: 0.6;
  }

  .tab-close-button.disabled {
    opacity: 0.2;
    cursor: not-allowed;
  }

  .tab-close-button:hover {
    background: rgba(244, 67, 54, 0.1);
    color: var(--error-color);
    opacity: 1;
    transform: translateY(-50%) scale(1.1);
  }

  .tab-container:hover .tab-close-button {
    opacity: 1;
  }

  .tab-add-button {
    padding: var(--spacing-sm) var(--spacing-md);
    border: none;
    border-radius: var(--radius-medium);
    font-size: 0.9rem;
    font-weight: 500;
    cursor: pointer;
    background: linear-gradient(
      135deg,
      var(--secondary-color),
      var(--secondary-light)
    );
    color: var(--text-primary);
    transition: all var(--transition-standard);
    box-shadow: var(--shadow-1);
    backdrop-filter: blur(10px);
    outline: none; /* 移除外框 */
  }

  .tab-add-button:hover {
    transform: translateY(-2px);
    box-shadow: var(--shadow-2);
    background: linear-gradient(
      135deg,
      var(--secondary-dark),
      var(--secondary-color)
    );
  }

  .tab-add-button:focus,
  .tab-add-button:active {
    outline: none; /* 確保在 focus 和 active 狀態下也沒有外框 */
    box-shadow: var(--shadow-1); /* 維持點擊時的陰影效果，避免預設 focus 樣式 */
  }

  /* 功能列樣式 */
  .function-bar {
    background: linear-gradient(
      135deg,
      var(--primary-color),
      var(--primary-dark)
    );
    padding: var(--spacing-md) var(--spacing-lg);
    box-shadow: var(--shadow-3);
    position: relative;
    z-index: 9;
  }

  .function-bar::before {
    content: "";
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 1px;
    background: linear-gradient(
      90deg,
      transparent,
      rgba(255, 255, 255, 0.5),
      transparent
    );
  }

  .function-buttons {
    display: flex;
    gap: var(--spacing-sm);
    align-items: center;
    flex-wrap: wrap;
  }

  .function-button {
    padding: var(--spacing-sm) var(--spacing-md);
    border: none;
    border-radius: var(--radius-medium);
    font-size: 0.9rem;
    font-weight: 500;
    cursor: pointer;
    background: rgba(255, 255, 255, 0.15);
    color: var(--text-on-primary);
    transition: all var(--transition-fast);
    backdrop-filter: blur(10px);
    border: 1px solid rgba(255, 255, 255, 0.2);
    position: relative;
    overflow: hidden;
  }

  .function-button::before {
    content: "";
    position: absolute;
    top: 0;
    left: -100%;
    width: 100%;
    height: 100%;
    background: linear-gradient(
      90deg,
      transparent,
      rgba(255, 255, 255, 0.3),
      transparent
    );
    transition: left 0.5s ease;
  }

  .function-button:hover {
    background: rgba(255, 255, 255, 0.25);
    transform: translateY(-1px);
    box-shadow: var(--shadow-2);
  }

  .function-button:hover::before {
    left: 100%;
  }

  /* 計算欄輸入區域 */
  .column-input {
    margin-top: var(--spacing-sm);
    padding: var(--spacing-sm);
    background: rgba(255, 255, 255, 0.1);
    border-radius: var(--radius-medium);
    backdrop-filter: blur(10px);
    border: 1px solid rgba(255, 255, 255, 0.2);
  }

  .input-row {
    display: flex;
    align-items: center;
    gap: var(--spacing-sm);
    flex-wrap: wrap;
  }

  .fx-label {
    color: rgba(255, 255, 255, 0.9);
    font-size: 0.9rem;
    font-weight: 600;
    min-width: 20px;
  }

  .column-name-input,
  .formula-input {
    padding: var(--spacing-xs) var(--spacing-sm);
    border: 1px solid rgba(255, 255, 255, 0.3);
    border-radius: var(--radius-small);
    background: rgba(255, 255, 255, 0.2);
    color: white;
    font-size: 0.9rem;
    backdrop-filter: blur(10px);
    transition: all var(--transition-fast);
  }

  .column-name-input {
    width: 120px;
    min-width: 120px;
  }

  .formula-input {
    flex: 1;
    min-width: 200px;
  }

  .column-name-input:focus,
  .formula-input:focus {
    border-color: rgba(255, 255, 255, 0.6);
    background: rgba(255, 255, 255, 0.3);
    box-shadow: 0 0 0 2px rgba(255, 255, 255, 0.2);
  }

  .column-name-input::placeholder,
  .formula-input::placeholder {
    color: rgba(255, 255, 255, 0.6);
  }

  .equals {
    color: rgba(255, 255, 255, 0.9);
    font-weight: 600;
    font-size: 0.9rem;
  }

  .confirm-button,
  .cancel-button {
    padding: var(--spacing-xs) var(--spacing-sm);
    border: none;
    border-radius: var(--radius-small);
    font-size: 0.8rem;
    font-weight: 500;
    cursor: pointer;
    transition: all var(--transition-fast);
    min-width: 50px;
  }

  .confirm-button {
    background: linear-gradient(135deg, var(--success-color), #66bb6a);
    color: white;
  }

  .cancel-button {
    background: linear-gradient(135deg, var(--error-color), #ef5350);
    color: white;
  }

  .confirm-button:hover,
  .cancel-button:hover {
    transform: translateY(-1px);
    box-shadow: var(--shadow-2);
  }

  .error-message {
    /* Existing styles */
    background-color: #fff5f5; /* Very light red background */
    color: #e53e3e; /* Softer red text */
    border: 1px solid #feb2b2; /* Light red border */
    padding: 0.5em;
    border-radius: 4px;
    margin-top: 0.5em;
    font-size: 0.9em;
  }

  /* 主要內容區域 */
  .main-content {
    flex: 1;
    display: flex;
    overflow: hidden;
    position: relative;
    z-index: 1;
    gap: var(--spacing-sm);
    padding: var(--spacing-sm);
  }

  .table-area {
    flex: 3;
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(20px);
    border-radius: var(--radius-large);
    box-shadow: var(--shadow-2);
    border: 1px solid rgba(255, 255, 255, 0.2);
    overflow: hidden;
    position: relative;
  }

  .table-area::before {
    content: "";
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 1px;
    background: linear-gradient(
      90deg,
      transparent,
      rgba(255, 255, 255, 0.6),
      transparent
    );
  }

  .table-placeholder {
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-direction: column;
    color: var(--text-secondary);
    font-size: 1.1rem;
    font-weight: 500;
  }

  .table-placeholder::before {
    content: "📊";
    font-size: 3rem;
    margin-bottom: var(--spacing-md);
    opacity: 0.5;
  }

  /* 右側資訊區域 */
  .info-area {
    flex: 1;
    background: rgba(255, 255, 255, 0.9);
    backdrop-filter: blur(20px);
    border-radius: var(--radius-large);
    box-shadow: var(--shadow-2);
    border: 1px solid rgba(255, 255, 255, 0.2);
    overflow: hidden;
    position: relative;
  }

  .info-area::before {
    content: "";
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 1px;
    background: linear-gradient(
      90deg,
      transparent,
      rgba(255, 255, 255, 0.6),
      transparent
    );
  }

  .info-header {
    padding: var(--spacing-lg);
    background: linear-gradient(
      135deg,
      rgba(25, 118, 210, 0.1),
      rgba(3, 218, 198, 0.1)
    );
    backdrop-filter: blur(10px);
    border-bottom: 1px solid rgba(255, 255, 255, 0.2);
  }

  .info-header h3 {
    margin: 0;
    font-size: 1.1rem;
    font-weight: 600;
    color: var(--primary-color);
    display: flex;
    align-items: center;
    gap: var(--spacing-sm);
  }

  .info-header h3::before {
    content: "📈";
    font-size: 1.2rem;
  }

  .stats-content {
    padding: var(--spacing-lg);
  }

  .stat-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: var(--spacing-sm) 0;
    border-bottom: 1px solid rgba(0, 0, 0, 0.06);
    transition: all var(--transition-fast);
  }

  .stat-item:hover {
    background: rgba(25, 118, 210, 0.05);
    margin: 0 calc(-1 * var(--spacing-sm));
    padding-left: var(--spacing-sm);
    padding-right: var(--spacing-sm);
    border-radius: var(--radius-small);
  }

  .stat-item:last-child {
    border-bottom: none;
  }

  .stat-label {
    font-size: 0.9rem;
    color: var(--text-secondary);
    font-weight: 500;
  }

  .stat-value {
    font-size: 0.9rem;
    font-weight: 600;
    color: var(--primary-color);
    padding: var(--spacing-xs) var(--spacing-sm);
    background: rgba(25, 118, 210, 0.1);
    border-radius: var(--radius-small);
    min-width: 40px;
    text-align: center;
  }

  /* 底部工具列 */
  .bottom-toolbar {
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(20px);
    padding: var(--spacing-sm);
    display: flex;
    gap: var(--spacing-xs);
    box-shadow: var(--shadow-2);
    border-top: 1px solid rgba(255, 255, 255, 0.2);
    position: relative;
    z-index: 10;
  }

  .bottom-toolbar::before {
    content: "";
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 1px;
    background: linear-gradient(
      90deg,
      transparent,
      rgba(255, 255, 255, 0.6),
      transparent
    );
  }

  .toolbar-button {
    flex: 1;
    padding: var(--spacing-sm) var(--spacing-md);
    border: none;
    border-radius: var(--radius-medium);
    font-size: 0.9rem;
    font-weight: 500;
    cursor: pointer;
    transition: all var(--transition-standard);
    backdrop-filter: blur(10px);
    border: 1px solid rgba(255, 255, 255, 0.3);
    position: relative;
    overflow: hidden;
  }

  .toolbar-button::before {
    content: "";
    position: absolute;
    top: 0;
    left: -100%;
    width: 100%;
    height: 100%;
    background: linear-gradient(
      90deg,
      transparent,
      rgba(255, 255, 255, 0.3),
      transparent
    );
    transition: left 0.5s ease;
  }

  .toolbar-button:hover::before {
    left: 100%;
  }

  .open-button {
    background: linear-gradient(
      135deg,
      rgba(76, 175, 80, 0.1),
      rgba(129, 199, 132, 0.1)
    );
    color: var(--success-color);
    border-color: rgba(76, 175, 80, 0.3);
  }

  .save-button {
    background: linear-gradient(
      135deg,
      rgba(33, 150, 243, 0.1),
      rgba(100, 181, 246, 0.1)
    );
    color: var(--info-color);
    border-color: rgba(33, 150, 243, 0.3);
  }

  .export-button {
    background: linear-gradient(
      135deg,
      rgba(255, 152, 0, 0.1),
      rgba(255, 183, 77, 0.1)
    );
    color: var(--warning-color);
    border-color: rgba(255, 152, 0, 0.3);
  }

  .toolbar-button:hover {
    transform: translateY(-2px);
    box-shadow: var(--shadow-2);
  }

  .open-button:hover {
    background: linear-gradient(
      135deg,
      rgba(76, 175, 80, 0.15),
      rgba(129, 199, 132, 0.15)
    );
  }

  .save-button:hover {
    background: linear-gradient(
      135deg,
      rgba(33, 150, 243, 0.15),
      rgba(100, 181, 246, 0.15)
    );
  }

  .export-button:hover {
    background: linear-gradient(
      135deg,
      rgba(255, 152, 0, 0.15),
      rgba(255, 183, 77, 0.15)
    );
  }

  /* 標籤名稱編輯樣式 */
  .tab-name-input {
    padding: var(--spacing-xs) var(--spacing-sm);
    border: 1px solid rgba(25, 118, 210, 0.3);
    border-radius: var(--radius-small);
    background: rgba(255, 255, 255, 0.95);
    color: var(--text-primary);
    font-size: 0.9rem;
    font-weight: 500;
    width: 120px;
    min-width: 120px;
    transition: all var(--transition-fast);
    backdrop-filter: blur(10px);
  }

  .tab-name-input:focus {
    border-color: var(--primary-color);
    box-shadow: 0 0 0 2px rgba(25, 118, 210, 0.2);
    background: rgba(255, 255, 255, 1);
  }

  .tab-name-input::placeholder {
    color: var(--text-hint);
  }

  /* 響應式設計 */
  @media (max-width: 1024px) {
    .main-content {
      flex-direction: column;
    }

    .info-area {
      flex: none;
      height: 200px;
    }

    .function-buttons {
      justify-content: center;
    }

    .input-row {
      flex-direction: column;
      align-items: stretch;
    }

    .formula-input {
      min-width: auto;
    }
  }

  @media (max-width: 768px) {
    .tab-row {
      overflow-x: auto;
      flex-wrap: nowrap;
    }

    .function-buttons {
      flex-wrap: wrap;
      gap: var(--spacing-xs);
    }

    .function-button {
      padding: var(--spacing-xs) var(--spacing-sm);
      font-size: 0.8rem;
    }

    .bottom-toolbar {
      flex-direction: column;
      gap: var(--spacing-xs);
    }

    .main-content {
      padding: var(--spacing-xs);
      gap: var(--spacing-xs);
    }
  }
</style>
