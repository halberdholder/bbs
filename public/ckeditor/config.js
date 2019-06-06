/**
 * @license Copyright (c) 2003-2018, CKSource - Frederico Knabben. All rights reserved.
 * For licensing, see https://ckeditor.com/legal/ckeditor-oss-license
 */

CKEDITOR.editorConfig = function( config ) {
	// Define changes to default configuration here.
	// For complete reference see:
	// http://docs.ckeditor.com/#!/api/CKEDITOR.config

	// The toolbar groups arrangement, optimized for two toolbar rows.
	//config.toolbarGroups = [
	//	{ name: 'clipboard',   groups: [ 'clipboard', 'undo' ] },
	//	{ name: 'editing',     groups: [ 'find', 'selection', 'spellchecker' ] },
	//	{ name: 'links' },
	//	{ name: 'insert' },
	//	{ name: 'forms' },
	//	{ name: 'tools' },
	//	{ name: 'document',	   groups: [ 'mode', 'document', 'doctools' ] },
	//	{ name: 'others' },
	//	'/',
	//	{ name: 'basicstyles', groups: [ 'basicstyles', 'cleanup' ] },
	//	{ name: 'paragraph',   groups: [ 'list', 'indent', 'blocks', 'align', 'bidi' ] },
	//	{ name: 'styles' },
	//	{ name: 'colors' },
	//	{ name: 'about' }
	//];

    //config.extraPlugins = 'codesnippet';

    //config.codeSnippet_languages = {
    //    javascript: 'JavaScript',
    //    php: 'PHP'
    //};
    
    config.toolbar = 'Basic';

    config.enterMode = CKEDITOR.ENTER_BR; 
    config.shiftEnterMode = CKEDITOR.ENTER_P;

	// Remove some buttons provided by the standard plugins, which are
	// not needed in the Standard(s) toolbar.
    config.removePlugins = 'elementspath';
	config.removeButtons = 'Underline,Subscript,Superscript';

	// Set the most common block elements.
	//config.format_tags = 'p;h1;h2;h3;pre';

	// Simplify the dialog windows.
	config.removeDialogTabs = 'image:advanced;link:advanced';

    config.extraPlugins = 'uploadimage';
    config.filebrowserUploadUrl = '/ckeditorUpload?type=File';  
    config.filebrowserImageUploadUrl = "/ckeditorUpload?type=image";
    config.filebrowserUploadMethod = "form";

    config.uiColor = '#f1e4db';
    config.allowedContent = false; // 是否允许使用源码模式进行编辑
    config.forcePasteAsPlainText = true; // 是否强制复制过来的文字去除格式
    config.keystrokes = [
      [CKEDITOR.CTRL + 86 /* V */, 'paste']
    ];
    config.blockedKeystrokes = [
      CKEDITOR.CTRL + 86
    ];

};
