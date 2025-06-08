<script lang="ts">
  import DataTable from "./components/DataTable.svelte";
  import Alert from "./components/Alert.svelte";
  import Confirm from "./components/Confirm.svelte";
  import Input from "./components/Input.svelte";
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

  // 計算欄輸入狀態
  let showColumnInput = false;
  let columnFormulaValue = "";
  let columnNameValue = "";
  let errorMessage = "";
  let showError = false;

  // 統計數據
  let currentStats = {
    總行數: "0",
    總欄數: "0",
    總儲存格: "0",
    數值欄數: "0",
  };

  // 組件掛載時執行
  onMount(async () => {
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
    const newTabName = `Tab ${tabs.length + 1}`;
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
          title: "創建失敗",
          message: "無法創建資料表",
          type: "error",
        });
        return;
      }
    }

    const columnName = await showInput({
      title: "新增欄位",
      message: "請輸入新欄位名稱:",
      placeholder: "欄位名稱",
      defaultValue: `新欄位 ${currentStats["總欄數"] ? parseInt(currentStats["總欄數"]) + 1 : 1}`,
      confirmText: "新增",
      cancelText: "取消",
    });
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

  function addCalculatedColumn() {
    showColumnInput = true;
  }

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
        cancelAddColumn();
      } else {
        showError = true;
        errorMessage = "添加計算欄失敗";
      }
    } else {
      showError = true;
      errorMessage = "請輸入欄位名稱與 CCL 表達式";
    }
  }

  function cancelAddColumn() {
    showColumnInput = false;
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

  // 接收統計數據更新
  function handleStatsUpdate(event: CustomEvent) {
    currentStats = event.detail;
  }

  // 底部工具列操作
  async function openFile() {
    // 簡單的文件選擇邏輯，可以擴展為文件對話框
    const input = await showInput({
      title: "開啟檔案",
      message: "請輸入檔案路徑:",
      placeholder: "檔案路徑",
      defaultValue: filePath || "test-data.json",
      confirmText: "開啟",
      cancelText: "取消",
    });
    if (input) {
      filePath = input;
      await handleLoadTable();
    }
  }

  async function saveFile() {
    if (!isTableLoaded) {
      await showAlert({
        title: "無法保存",
        message: "沒有已載入的資料表可保存",
        type: "warning",
      });
      return;
    }

    // 如果沒有指定路徑，提示用戶輸入
    if (!filePath) {
      const input = await showInput({
        title: "儲存檔案",
        message: "請輸入保存路徑:",
        placeholder: "檔案路徑",
        defaultValue: "saved-data.json",
        confirmText: "儲存",
        cancelText: "取消",
      });
      if (input) {
        filePath = input;
      } else {
        return;
      }
    }

    await handleSaveTable();
  }

  async function exportFile() {
    // TODO: 實現匯出功能，可以支持 CSV, Excel 等格式
    console.log("匯出檔案", tabs[currentTabIndex]?.id);
    await showAlert({
      title: "功能開發中",
      message: "匯出功能尚未實現",
      type: "info",
    });
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
    if (!filePath) {
      await showAlert({
        title: "載入錯誤",
        message: "請輸入檔案路徑",
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
</script>

<main>
  <!-- 標籤列 -->
  <div class="tab-bar">
    <div class="tab-row">
      {#each tabs as tab, index}
        <div class="tab-container">
          <button
            class="tab-button"
            class:tab-active={tab.isActive}
            on:click={() => switchTab(index)}
          >
            {tab.name}
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
      <button class="tab-add-button" on:click={addNewTab}>+</button>
    </div>
  </div>

  <!-- 功能列 -->
  <div class="function-bar">
    <div class="function-buttons">
      <button class="function-button" on:click={addColumn}>新增欄</button>
      <button class="function-button" on:click={addRow}>新增列</button>
      <button class="function-button" on:click={addCalculatedColumn}
        >新增計算欄</button
      >
    </div>

    <!-- 計算欄輸入區域 -->
    {#if showColumnInput}
      <div class="column-input">
        <div class="input-row">
          <span class="fx-label">fx</span>
          <input
            type="text"
            class="column-name-input"
            placeholder="名稱"
            bind:value={columnNameValue}
          />
          <span class="equals">=</span>
          <input
            type="text"
            class="formula-input"
            placeholder="CCL 表達式"
            bind:value={columnFormulaValue}
          />
          <button class="confirm-button" on:click={confirmAddColumn}>✓</button>
          <button class="cancel-button" on:click={cancelAddColumn}>✕</button>
        </div>
        {#if showError}
          <div class="error-message">{errorMessage}</div>
        {/if}
      </div>
    {/if}
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
          <p>資料表為空，請新增資料</p>
        </div>
      {/if}
    </div>

    <!-- 右側資訊區域 -->
    <div class="info-area">
      <div class="info-header">
        <h3>基本統計</h3>
      </div>
      <div class="stats-content">
        <div class="stat-item">
          <span class="stat-label">行數:</span>
          <span class="stat-value">{currentStats["總行數"]}</span>
        </div>
        <div class="stat-item">
          <span class="stat-label">列數:</span>
          <span class="stat-value">{currentStats["總欄數"]}</span>
        </div>
        <div class="stat-item">
          <span class="stat-label">總儲存格:</span>
          <span class="stat-value">{currentStats["總儲存格"]}</span>
        </div>
        <div class="stat-item">
          <span class="stat-label">數值欄數:</span>
          <span class="stat-value">{currentStats["數值欄數"]}</span>
        </div>
      </div>
    </div>
  </div>

  <!-- 底部工具列 -->
  <div class="bottom-toolbar">
    <button class="toolbar-button open-button" on:click={openFile}>開啟</button>
    <button class="toolbar-button save-button" on:click={saveFile}>存檔</button>
    <button class="toolbar-button export-button" on:click={exportFile}
      >匯出</button
    >
    <button class="toolbar-button settings-button" on:click={openSettings}
      >設定</button
    >
  </div>
</main>

<!-- 對話框組件 -->
<Alert
  visible={$alertStore.visible}
  options={$alertStore.options}
  on:close={(e) => closeAlert()}
/>

<Confirm
  visible={$confirmStore.visible}
  options={$confirmStore.options}
  on:close={(e) => closeConfirm(e.detail.result)}
/>

<Input
  visible={$inputStore.visible}
  options={$inputStore.options}
  on:close={(e) => closeInput(e.detail.result)}
/>

<style>
  :global(body) {
    margin: 0;
    padding: 0;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen,
      Ubuntu, Cantarell, "Open Sans", "Helvetica Neue", sans-serif;
  }

  main {
    display: flex;
    flex-direction: column;
    height: 100vh;
    width: 100%;
    background-color: #f5f8ff;
  }

  /* 標籤列樣式 */
  .tab-bar {
    background-color: #ffffff;
    border-bottom: 1px solid #e0e0e0;
    padding: 4px 8px 0 8px;
  }

  .tab-row {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;
    align-items: center;
  }

  .tab-container {
    position: relative;
    display: flex;
    align-items: center;
  }

  .tab-button {
    padding: 8px 16px;
    border: none;
    border-radius: 4px 4px 0 0;
    font-size: 14px;
    cursor: pointer;
    transition: all 0.2s;
    background-color: rgb(225, 235, 250); /* 淡藍色背景 - 未選中 */
    color: rgb(0, 90, 180); /* 藍色文字 */
    margin-bottom: -1px;
    padding-right: 32px; /* 為刪除按鈕留出空間 */
    position: relative;
    overflow: visible;
  }

  .tab-button.tab-active {
    background-color: rgb(235, 250, 235); /* 淡綠色背景 - 選中 */
    color: rgb(0, 0, 0); /* 黑色文字 */
    border-bottom: 2px solid rgb(235, 250, 235);
  }

  .tab-button:hover {
    opacity: 0.8;
  }

  .tab-close-button {
    position: absolute;
    right: 4px;
    top: 50%;
    transform: translateY(-50%);
    width: 20px;
    height: 20px;
    border: none;
    border-radius: 50%;
    background-color: rgba(0, 0, 0, 0.1);
    color: rgb(0, 90, 180);
    font-size: 16px;
    font-weight: bold;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s;
    line-height: 1;
    padding: 0;
    z-index: 10;
  }

  .tab-close-button.disabled {
    opacity: 0.3;
    cursor: not-allowed;
  }

  .tab-close-button:hover {
    background-color: rgba(255, 0, 0, 0.2);
    color: rgb(200, 0, 0);
  }

  .tab-container:hover .tab-close-button {
    background-color: rgba(0, 0, 0, 0.15);
  }

  .tab-add-button {
    padding: 8px 12px;
    border: none;
    border-radius: 4px;
    font-size: 14px;
    cursor: pointer;
    background-color: rgb(225, 245, 254); /* 淡藍色背景 */
    color: rgb(33, 150, 243); /* 藍色文字 */
    margin-bottom: -1px;
  }

  .tab-add-button:hover {
    opacity: 0.8;
  }

  /* 功能列樣式 */
  .function-bar {
    background-color: rgb(63, 81, 181); /* Material Design Indigo 500 */
    padding: 8px 16px;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  }

  .function-buttons {
    display: flex;
    gap: 8px;
    align-items: center;
  }

  .function-button {
    padding: 8px 16px;
    border: none;
    border-radius: 4px;
    font-size: 14px;
    cursor: pointer;
    background-color: rgba(255, 255, 255, 0.15); /* 半透明白色背景 */
    color: rgb(255, 255, 255); /* 白色文字 */
    transition: all 0.2s;
  }

  .function-button:hover {
    background-color: rgba(255, 255, 255, 0.25);
  }

  /* 計算欄輸入區域 */
  .column-input {
    margin-top: 8px;
    padding: 8px;
    background-color: rgba(255, 255, 255, 0.1);
    border-radius: 4px;
  }

  .input-row {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .fx-label {
    color: rgb(200, 200, 200);
    font-size: 14px;
    font-weight: bold;
  }

  .column-name-input {
    width: 100px;
    padding: 4px 8px;
    border: 1px solid rgba(255, 255, 255, 0.3);
    border-radius: 4px;
    background-color: rgba(255, 255, 255, 0.2);
    color: white;
    font-size: 14px;
  }

  .column-name-input::placeholder {
    color: rgba(255, 255, 255, 0.7);
  }

  .equals {
    color: white;
    font-weight: bold;
  }

  .formula-input {
    flex: 1;
    padding: 4px 8px;
    border: 1px solid rgba(255, 255, 255, 0.3);
    border-radius: 4px;
    background-color: rgba(255, 255, 255, 0.2);
    color: white;
    font-size: 14px;
  }

  .formula-input::placeholder {
    color: rgba(255, 255, 255, 0.7);
  }

  .confirm-button {
    padding: 4px 8px;
    border: none;
    border-radius: 4px;
    background-color: rgb(0, 150, 0);
    color: white;
    cursor: pointer;
    font-size: 12px;
  }

  .cancel-button {
    padding: 4px 8px;
    border: none;
    border-radius: 4px;
    background-color: rgb(150, 0, 0);
    color: white;
    cursor: pointer;
    font-size: 12px;
  }

  .error-message {
    margin-top: 4px;
    padding: 4px 8px;
    color: rgb(255, 200, 200);
    font-size: 12px;
  }

  /* 主要內容區域 */
  .main-content {
    flex: 1;
    display: flex;
    overflow: hidden;
  }

  .table-area {
    flex: 3;
    background-color: white;
    overflow: auto;
  }

  .table-placeholder {
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #757575;
    font-size: 16px;
  }

  /* 右側資訊區域 */
  .info-area {
    flex: 1;
    background-color: rgb(245, 248, 255); /* 淡藍灰色 */
    border-left: 1px solid #e0e0e0;
    overflow: auto;
  }

  .info-header {
    padding: 16px;
    border-bottom: 1px solid #e0e0e0;
  }

  .info-header h3 {
    margin: 0;
    font-size: 16px;
    font-weight: bold;
    color: rgb(0, 90, 180); /* 藍色 */
  }

  .stats-content {
    padding: 16px;
  }

  .stat-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 8px 0;
    border-bottom: 1px solid rgba(0, 0, 0, 0.1);
  }

  .stat-label {
    font-size: 14px;
    color: #666;
  }

  .stat-value {
    font-size: 14px;
    font-weight: bold;
    color: #333;
  }

  /* 底部工具列 */
  .bottom-toolbar {
    background-color: rgb(180, 220, 255); /* 淡藍紫色背景 */
    padding: 6px;
    display: flex;
    box-shadow: 0 -1px 3px rgba(0, 0, 0, 0.1);
    border-top: 1px solid rgb(220, 225, 230);
  }

  .toolbar-button {
    flex: 1;
    padding: 8px 16px;
    border: none;
    border-radius: 0;
    font-size: 14px;
    cursor: pointer;
    margin: 0 3px;
    transition: all 0.2s;
  }

  .open-button {
    background-color: rgb(240, 253, 244); /* 淡綠色背景 */
    color: rgb(34, 197, 94); /* 綠色文字 */
  }

  .save-button {
    background-color: rgb(239, 246, 255); /* 淡藍色背景 */
    color: rgb(59, 130, 246); /* 藍色文字 */
  }

  .export-button {
    background-color: rgb(255, 251, 235); /* 淡橙色背景 */
    color: rgb(245, 158, 11); /* 橙色文字 */
  }

  .settings-button {
    background-color: rgb(249, 250, 251); /* 淡灰色背景 */
    color: rgb(107, 114, 128); /* 灰色文字 */
  }

  .toolbar-button:hover {
    opacity: 0.8;
  }
</style>
