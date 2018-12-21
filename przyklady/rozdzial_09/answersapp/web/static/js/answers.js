$(function(){

	$.fn.showModal = function(){
		$(this).modal('setting', {blurring:true}).modal('show')
	}

	// object pobiera obiekt opisujący wszystkie wartości wejściowe z formularza.
	$.fn.object = function(){
		var obj = {}
		$(this).find('input,textarea').each(function(){
			var $this = $(this),
				name = $this.attr('name'),
				value = $this.val()
			if (name && value)
				obj[name] = value
		})
		return obj
	}

	// api to opakowanie dla $.ajax realizujące często wykonywane czynności
	$.api = function(options){
		if (options.form) {
			options.type = options.form.attr('method') || 'get'
			options.url = options.form.attr('action') || ""
			options.data = JSON.stringify(options.form.object())
			options.dataType = 'json'
			options.contentType = 'application/json'
		}
		options._error = options.error
		options.error = function(response){
			console.warn(response)
			var message = response.responseText || response.statusText || "Wystąpił nieznany błąd."
			if (response.responseJSON && response.responseJSON.error) {
				message = response.responseJSON.error
			}
			options._error = options._error || function(message){
				// global error handler
				var errEl = $('.ui.global.error.message')
				if (errEl.length == 0) {
					errEl = $("<div>", {class:"ui global error message container"}).insertAfter($(".topnav"))
				}
				errEl.text(message)
			}
			options._error(message, response)
		}
		$.ajax(options)
	}

	// inject ustawia tekst lub zawartość elementów stron zapisjując w nich 
	// dane przekazane jako argumenty wywołania.
	$.fn.inject = function(data){
		var $this = $(this)
		for (var k in data) {
			if (!data.hasOwnProperty(k)) continue
			$this.find('[data-field="'+k+'"]').each(function(){
				var $that = $(this)
				switch ($that[0].tagName) {
					case "INPUT":
						$that.val(data[k])
						break
					default:
						$that.text(data[k])
				}
			})
		}
	}

	// data-trigger-modal - zdarzenia
	$('[data-trigger-modal]').on('click', function(e){
		e.preventDefault()
		var modal = $(this).attr('data-trigger-modal')
		$('[data-modal="'+modal+'"]').showModal()
	})

	// ponownie instaluje wszystkie odroczone skrypty, które 
	// mogły już zostać odłączone do strony.
	$('script[type="deferred"]')
		.remove()
		.attr('type', 'text/javascript')
		.appendTo($('body'))

})